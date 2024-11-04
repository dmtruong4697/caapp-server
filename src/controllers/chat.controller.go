package controllers

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	responce_models "caapp-server/src/models/responce_models"
	utils "caapp-server/src/utils"
	"log"
	"net/http"
	"strconv"
	"time"

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

var channels = make(map[uint64]map[*websocket.Conn]bool)
var broadcast = make(chan responce_models.GetChannelChatHistoryItem)

func HandleConnections(c *gin.Context) {
	channelIDStr := c.Query("channel_id")

	channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer ws.Close()

	// Khởi tạo channel nếu chưa tồn tại
	if channels[channelID] == nil {
		channels[channelID] = make(map[*websocket.Conn]bool)
	}
	channels[channelID][ws] = true

	// Lắng nghe tin nhắn từ client và gửi vào kênh broadcast
	for {
		var msg db_models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(channels[channelID], ws)
			if len(channels[channelID]) == 0 {
				delete(channels, channelID)
			}
			break
		}

		msg.ChannelID = uint(channelID)
		msg.CreateAt = time.Now()
		msg.IsEdited = false
		msg.LastUpdate = time.Now()

		if err := database.DB.Create(&msg).Error; err != nil {
			// c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_"})
			// return
		}

		var msgDetail responce_models.GetChannelChatHistoryItem
		msgDetail.Message = msg
		msgDetail.Sender = utils.GetProfileInfo(msg.SenderID)

		// Gửi tin nhắn vào kênh broadcast
		broadcast <- msgDetail
	}
}

func HandleMessages() {
	for {
		// Lấy tin nhắn từ kênh broadcast
		msgDetail := <-broadcast
		channelID64 := uint64(msgDetail.Message.ChannelID)

		// Gửi tin nhắn tới tất cả các client trong channel tương ứng
		for client := range channels[channelID64] {
			err := client.WriteJSON(msgDetail)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(channels[channelID64], client)
				if len(channels[channelID64]) == 0 {
					delete(channels, channelID64)
				}
			}
		}
	}
}
