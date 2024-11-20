package ws

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	responce_models "caapp-server/src/models/responce_models"
	wsmodels "caapp-server/src/models/ws_models"
	"caapp-server/src/utils"
	"caapp-server/src/utils/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var channels = make(map[string]map[*websocket.Conn]bool)
var broadcast = make(chan wsmodels.WSChannelListItemForHandleMessage)

func SendToBroadcast(msg wsmodels.WSChannelListItemForHandleMessage) {
	broadcast <- msg
}

func HandleChannelListConnections(c *gin.Context) {
	userID := c.Query("user_id")

	channelID := userID + "-ChatListChannel"

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer ws.Close()

	if channels[channelID] == nil {
		channels[channelID] = make(map[*websocket.Conn]bool)
	}
	channels[channelID][ws] = true

	for {
		var msg wsmodels.WSChannelListItemForListening
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(channels[channelID], ws)
			if len(channels[channelID]) == 0 {
				delete(channels, channelID)
			}
			break
		}

		var msgForBroadcast wsmodels.WSChannelListItemForBroadcast
		msgForBroadcast.Channel = utils.GetChannelInfo(msg.ChannelID)

		var lastMessage db_models.Message
		if err := database.DB.Where("id = ?", msgForBroadcast.Channel.LastMessageID).First(&lastMessage).Error; err != nil {

		}
		msgForBroadcast.LastMessage = lastMessage
		msgForBroadcast.LastMessageSender = utils.GetUserInfo(helper.StringToUInt(userID), msgForBroadcast.LastMessage.SenderID)

		// query all channel member info
		// query := `
		// 	SELECT *
		// 	FROM channel_members
		// 	WHERE channel_id = ?
		// 	AND user_id != ?;
		// `

		var channelMembers []db_models.ChannelMember
		err = database.DB.Where("channel_id = ?", msg.ChannelID).Find(&channelMembers).Error
		if err != nil {
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channels"})
			// return
		}

		msgForBroadcast.Users = make([]responce_models.GetUserInfoResponce, len(channelMembers))
		for j := range channelMembers {
			msgForBroadcast.Users[j] = utils.GetUserInfo(helper.StringToUInt(userID), channelMembers[j].UserID)
		}

		var msgForHandleMessage wsmodels.WSChannelListItemForHandleMessage
		msgForHandleMessage.Data = msgForBroadcast
		msgForHandleMessage.UserIDForMakingChannelID = userID
		broadcast <- msgForHandleMessage
	}
}

func HandleChatListMessages() {
	for {
		msg := <-broadcast

		channelID := msg.UserIDForMakingChannelID + "-ChatListChannel"

		for client := range channels[channelID] {
			err := client.WriteJSON(msg.Data)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(channels[channelID], client)
				if len(channels[channelID]) == 0 {
					delete(channels, channelID)
				}
			}
		}
	}
}
