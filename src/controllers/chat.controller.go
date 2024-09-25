package controllers

import (
	db_models "caapp-server/src/models/db_models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
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

	for {
		var msg db_models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(channels[channelID], ws)
			break
		}

		for client := range channels[channelID] {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(channels[channelID], client)
			}
		}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast

		channelID64 := uint64(msg.ChannelID)

		for client := range channels[channelID64] {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(channels[channelID64], client)
			}
		}
	}
}
