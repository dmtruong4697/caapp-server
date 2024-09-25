package controllers

import (
	db_models "caapp-server/src/models/db_models"
	"log"
	"net/http"
	"strconv"

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
var broadcast = make(chan db_models.Message)

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

		// Thêm ChannelID cho message để xác định channel
		msg.ChannelID = uint(channelID)

		// Gửi tin nhắn vào kênh broadcast
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		// Lấy tin nhắn từ kênh broadcast
		msg := <-broadcast
		channelID64 := uint64(msg.ChannelID)

		// Gửi tin nhắn tới tất cả các client trong channel tương ứng
		for client := range channels[channelID64] {
			err := client.WriteJSON(msg)
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
