package rcroutines

import (
	"caapp-server/src/database"
	models "caapp-server/src/models/db_models"
	rcdbmodels "caapp-server/src/models/db_models/rc_db_models"
	"caapp-server/src/utils/helper"
	rcws "caapp-server/src/ws/rc_ws"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var queueMutex sync.Mutex

func PairUser() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		queueMutex.Lock()

		var allUserInQueue []rcdbmodels.GeneralQueueUser
		database.DB.Order("join_at ASC").Find(&allUserInQueue)

		// check user number in queue
		if len(allUserInQueue) < 2 {
			queueMutex.Unlock()
			continue
		}

		var userA, userB rcdbmodels.GeneralQueueUser
		firstUserInQueue := allUserInQueue[0]

		// finding user in queue
		for i := 1; i < len(allUserInQueue); i++ {
			if IsMatchingUser(firstUserInQueue, allUserInQueue[i]) {
				userA = firstUserInQueue
				userB = allUserInQueue[i]
				break
			}
		}

		// check user
		if userA.UserID == 0 || userB.UserID == 0 {
			queueMutex.Unlock()
			continue
		}

		// create new RC Channel
		newRCChannel := rcdbmodels.RCChannel{
			CreateAt: time.Now(),
		}
		if err := database.DB.Create(&newRCChannel).Error; err != nil {
			log.Printf("Failed to create RC Channel: %v", err)
			queueMutex.Unlock()
			continue
		}

		// create channel member
		channelMembers := []rcdbmodels.RCChannelMember{
			{UserID: userA.UserID, ChannelID: newRCChannel.ID, JoinAt: time.Now()},
			{UserID: userB.UserID, ChannelID: newRCChannel.ID, JoinAt: time.Now()},
		}
		if err := database.DB.Create(&channelMembers).Error; err != nil {
			log.Printf("Failed to create channel members: %v", err)
		}

		// update user's CurrentRCChannelID
		updateUserChannel(userA.UserID, newRCChannel.ID)
		updateUserChannel(userB.UserID, newRCChannel.ID)

		// send notification to ws queue channel
		notifyWebSocket(userA.UserID, newRCChannel)
		notifyWebSocket(userB.UserID, newRCChannel)

		// delete user from queue
		removeFromQueue(userA.UserID)
		removeFromQueue(userB.UserID)

		queueMutex.Unlock()
	}
}

func IsMatchingUser(userA rcdbmodels.GeneralQueueUser, userB rcdbmodels.GeneralQueueUser) bool {
	return userA.TargetGender == userB.Gender && userB.TargetGender == userA.Gender
}

func updateUserChannel(userID uint, channelID uint) {
	var dbUser models.User
	if err := database.DB.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		log.Printf("Failed to find user %d: %v", userID, err)
		return
	}

	dbUser.CurrentRCChannelID = channelID
	if err := database.DB.Save(&dbUser).Error; err != nil {
		log.Printf("Failed to update user %d RC Channel: %v", userID, err)
	}
}

func notifyWebSocket(userID uint, channel rcdbmodels.RCChannel) {
	channelID := helper.UIntToString(userID) + "-RCQueueChannel"
	queueChannels := rcws.GetQueueChannels()

	if clients, ok := queueChannels[channelID]; ok {
		for client := range clients {
			err := client.WriteJSON(channel)
			if err != nil {
				log.Printf("WebSocket error for user %d: %v", userID, err)
				client.Close()
				delete(clients, client)

				if len(clients) == 0 {
					delete(queueChannels, channelID)
				}
			}
		}
	}
}

func removeFromQueue(userID uint) {
	if err := database.DB.Where("user_id = ?", userID).Delete(&rcdbmodels.GeneralQueueUser{}).Error; err != nil {
		log.Printf("Failed to remove user %d from queue: %v", userID, err)
	}
}
