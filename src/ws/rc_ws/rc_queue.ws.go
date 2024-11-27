package rcws

import (
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

func HandleWaitingQueueConnections(c *gin.Context) {
	userID := c.Query("user_id")

	channelID := userID + "-RCQueueChannel"

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

	}
}
