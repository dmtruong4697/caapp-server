package rcws

import (
	"caapp-server/src/database"
	rcdbmodels "caapp-server/src/models/db_models/rc_db_models"
	rcwsmodels "caapp-server/src/models/ws_models/rc_ws_models"
	"caapp-server/src/utils/helper"
	"fmt"
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

var queueChannels = make(map[string]map[*websocket.Conn]bool)
var broadcast = make(chan rcwsmodels.MatchingResponse)

func SendToWaitingQueueBroadcast(msg rcwsmodels.MatchingResponse) {
	broadcast <- msg
}

func GetQueueChannels() map[string]map[*websocket.Conn]bool {
	return queueChannels
}

func HandleWaitingQueueConnections(c *gin.Context) {
	userID := c.Query("user_id")

	channelID := userID + "-RCQueueChannel"

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	defer func() {
		ws.Close()
		delete(queueChannels[channelID], ws)
		if len(queueChannels[channelID]) == 0 {
			delete(queueChannels, channelID)
		}

		var generalQueueUser rcdbmodels.GeneralQueueUser
		if err := database.DB.Where("user_id = ?", helper.StringToUInt(userID)).First(&generalQueueUser).Error; err == nil {
			if err := database.DB.Delete(&generalQueueUser).Error; err != nil {
				log.Printf("error deleting user from general queue: %v", err)
			}
		}
	}()

	if queueChannels[channelID] == nil {
		queueChannels[channelID] = make(map[*websocket.Conn]bool)
	}
	queueChannels[channelID][ws] = true

	for {
		var msg rcwsmodels.RCQueueRequest
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(queueChannels[channelID], ws)
			if len(queueChannels[channelID]) == 0 {
				delete(queueChannels, channelID)
			}
			break
		}

		fmt.Println(msg)

		if msg.RequestType == "0" {
			var generalQueueUser rcdbmodels.GeneralQueueUser
			generalQueueUser.UserID = msg.UserID
			generalQueueUser.Gender = msg.Gender
			generalQueueUser.TargetGender = msg.TargetGender

			if err := database.DB.Create(&generalQueueUser).Error; err != nil {
			}
		} else if msg.RequestType == "1" {
			var generalQueueUser rcdbmodels.GeneralQueueUser
			if err := database.DB.Where("user_id = ?", msg.UserID).First(&generalQueueUser).Error; err != nil {
			}

			// Xóa bản ghi
			if err := database.DB.Where("user_id = ?", generalQueueUser.UserID).Delete(&generalQueueUser).Error; err != nil {
			}
		}
	}
}

func HandleWaitingQueueMessages() {
	for {
		msg := <-broadcast

		channelID := helper.UIntToString(msg.UserID) + "-RCQueueChannel"

		for client := range queueChannels[channelID] {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(queueChannels[channelID], client)
				if len(queueChannels[channelID]) == 0 {
					delete(queueChannels, channelID)
				}
			}
		}
	}
}
