package rcroutines

import (
	"caapp-server/src/database"
	rcdbmodels "caapp-server/src/models/db_models/rc_db_models"
	"caapp-server/src/utils/helper"
	rcws "caapp-server/src/ws/rc_ws"
	"log"
	"net/http"
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

// var queueChannels = make(map[string]map[*websocket.Conn]bool)

func PairUser() {
	for {
		var userA rcdbmodels.GeneralQueueUser
		var userB rcdbmodels.GeneralQueueUser

		var allUserInQueue []rcdbmodels.GeneralQueueUser
		database.DB.Order("join_at ASC").Find(&allUserInQueue)

		if len(allUserInQueue) < 2 {
			time.Sleep(1 * time.Second)
			continue
		}

		firstUserInQueue := allUserInQueue[0]
		for i := 1; i < len(allUserInQueue); i++ {
			if IsMatchingUser(firstUserInQueue, allUserInQueue[i]) {
				userB = allUserInQueue[i]
				userA = firstUserInQueue
				break
			}
		}

		channelIDA := helper.UIntToString(userA.UserID) + "-RCQueueChannel"
		channelIDB := helper.UIntToString(userB.UserID) + "-RCQueueChannel"

		// tao moi channel
		newRCChannel := rcdbmodels.RCChannel{
			CreateAt: time.Now(),
		}

		if err := database.DB.Create(&newRCChannel).Error; err != nil {
			// handle loi tao moi channel
		}

		// tao moi channel member
		channelMember1 := rcdbmodels.RCChannelMember{
			UserID:    userA.UserID,
			ChannelID: newRCChannel.ID,
			JoinAt:    time.Now(),
		}

		channelMember2 := rcdbmodels.RCChannelMember{
			UserID:    userB.UserID,
			ChannelID: newRCChannel.ID,
			JoinAt:    time.Now(),
		}

		if err := database.DB.Create(&channelMember1).Error; err != nil {
		}
		if err := database.DB.Create(&channelMember2).Error; err != nil {
		}

		for client := range rcws.GetQueueChannels()[channelIDA] {
			err := client.WriteJSON(newRCChannel)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(rcws.GetQueueChannels()[channelIDA], client)
				if len(rcws.GetQueueChannels()[channelIDA]) == 0 {
					delete(rcws.GetQueueChannels(), channelIDA)
				}
			}
		}

		for client := range rcws.GetQueueChannels()[channelIDB] {
			err := client.WriteJSON(newRCChannel)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(rcws.GetQueueChannels()[channelIDB], client)
				if len(rcws.GetQueueChannels()[channelIDB]) == 0 {
					delete(rcws.GetQueueChannels(), channelIDB)
				}
			}
		}

		// xoa 2 user khoi queue
		if err := database.DB.Delete(&userA).Error; err != nil {
		}
		if err := database.DB.Delete(&userB).Error; err != nil {
		}

	}
}

func IsMatchingUser(userA rcdbmodels.GeneralQueueUser, userB rcdbmodels.GeneralQueueUser) bool {
	if (userA.TargetGender != userB.Gender) || (userB.TargetGender != userA.Gender) {
		return false
	}
	return true
}
