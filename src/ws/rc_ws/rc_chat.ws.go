package rcws

import (
	"caapp-server/src/database"
	rcdbmodels "caapp-server/src/models/db_models/rc_db_models"
	rcresponsemodel "caapp-server/src/models/responce_models/rc_response_model"
	"caapp-server/src/utils/helper"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var RCChatChannels = make(map[string]map[*websocket.Conn]bool)
var RCChatBroadcast = make(chan rcresponsemodel.GetRCChannelChatHistoryItem)

func GetRCChatChannels() map[string]map[*websocket.Conn]bool {
	return RCChatChannels
}

func HandleRCChatConnections(c *gin.Context) {
	channelID := c.Query("channel_id") + "-RCChatChannel"

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer ws.Close()

	if RCChatChannels[channelID] == nil {
		RCChatChannels[channelID] = make(map[*websocket.Conn]bool)
	}
	RCChatChannels[channelID][ws] = true

	for {
		var msg rcdbmodels.RCMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(RCChatChannels[channelID], ws)
			if len(RCChatChannels[channelID]) == 0 {
				delete(RCChatChannels, channelID)
			}
			break
		}

		fmt.Println(msg)

		msg.ChannelID = helper.StringToUInt(c.Query("channel_id"))
		msg.CreateAt = time.Now()
		msg.IsEdited = false
		msg.LastUpdate = time.Now()

		if err := database.DB.Create(&msg).Error; err != nil {
		}

		var msgDetail rcresponsemodel.GetRCChannelChatHistoryItem
		msgDetail.Message = msg

		RCChatBroadcast <- msgDetail
	}
}

func HandleRCChatMessages() {
	for {
		msg := <-RCChatBroadcast

		channelID := helper.UIntToString(msg.Message.ChannelID) + "-RCChatChannel"

		for client := range RCChatChannels[channelID] {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(RCChatChannels[channelID], client)
				if len(RCChatChannels[channelID]) == 0 {
					delete(RCChatChannels, channelID)
				}
			}
		}
	}
}
