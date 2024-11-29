package controllers

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	responce_models "caapp-server/src/models/responce_models"
	wsmodels "caapp-server/src/models/ws_models"
	utils "caapp-server/src/utils"
	"caapp-server/src/utils/helper"
	ws "caapp-server/src/ws"
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
var chatListChannels = make(map[string]map[*websocket.Conn]bool)
var broadcast = make(chan responce_models.GetChannelChatHistoryItem)

func SendToChannelList(userID string, msg wsmodels.WSChannelListItemForBroadcast) {
	msgForHandleMessage := wsmodels.WSChannelListItemForHandleMessage{
		Data:                     msg,
		UserIDForMakingChannelID: userID,
	}
	ws.SendToChannelListBroadcast(msgForHandleMessage)
}

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

	// ket noi toi channel chat chinh
	if channels[channelID] == nil {
		channels[channelID] = make(map[*websocket.Conn]bool)
	}
	channels[channelID][ws] = true
	//////////////////////////////////////////////

	// ket noi toi cac chat list channel
	var channelMembers []db_models.ChannelMember
	database.DB.Where(
		"(channel_id = ?)",
		channelID,
	).Find(&channelMembers)

	for i := range channelMembers {
		chatListChannelID := helper.UIntToString(channelMembers[i].UserID) + "-ChatListChannel"
		if chatListChannels[chatListChannelID] == nil {
			chatListChannels[chatListChannelID] = make(map[*websocket.Conn]bool)
		}
		chatListChannels[chatListChannelID][ws] = true
	}
	/////////

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

		var channel db_models.Channel
		if err := database.DB.Where("id = ?", msg.ChannelID).First(&channel).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error_code": ""})
			return
		}

		channel.LastMessageID = msg.ID
		if err := database.DB.Save(&channel).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": ""})
			return
		}

		var msgDetail responce_models.GetChannelChatHistoryItem
		msgDetail.Message = msg
		msgDetail.Sender = utils.GetProfileInfo(msg.SenderID)

		for _, member := range channelMembers {
			// chatListChannelID := helper.UIntToString(channelMembers[i].UserID) + "-ChatListChannel"
			// var msg wsmodels.WSChannelListItemForListening

			// msg.ChannelID = msgDetail.Message.ChannelID
			// msg.Message = msgDetail.Message

			var msgForBroadcast wsmodels.WSChannelListItemForBroadcast
			msgForBroadcast.Channel = utils.GetChannelInfo(msgDetail.Message.ChannelID)

			var lastMessage db_models.Message
			if err := database.DB.Where("id = ?", msgForBroadcast.Channel.LastMessageID).First(&lastMessage).Error; err != nil {

			}
			msgForBroadcast.LastMessage = lastMessage
			msgForBroadcast.LastMessageSender = utils.GetUserInfo(msg.SenderID, msgForBroadcast.LastMessage.SenderID)

			// query all channel member info
			// query := `
			// 	SELECT *
			// 	FROM channel_members
			// 	WHERE channel_id = ?
			// 	AND user_id != ?;
			// `

			var channelMembers []db_models.ChannelMember
			err := database.DB.Where("channel_id = ?", msgDetail.Message.ChannelID).Find(&channelMembers).Error
			if err != nil {
				// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channels"})
				// return
			}

			msgForBroadcast.Users = make([]responce_models.GetUserInfoResponce, len(channelMembers))
			for j := range channelMembers {
				msgForBroadcast.Users[j] = utils.GetUserInfo(member.UserID, channelMembers[j].UserID)
			}
			SendToChannelList(helper.UIntToString(member.UserID), msgForBroadcast)
		}

		broadcast <- msgDetail
	}
}

func HandleMessages() {
	for {
		msgDetail := <-broadcast
		channelID64 := uint64(msgDetail.Message.ChannelID)

		var channelMembers []db_models.ChannelMember
		database.DB.Where(
			"(channel_id = ?)",
			msgDetail.Message.ChannelID,
		).Find(&channelMembers)

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

		// for i := range channelMembers {
		// 	chatListChannelID := helper.UIntToString(channelMembers[i].UserID) + "-ChatListChannel"
		// 	// var msg wsmodels.WSChannelListItemForListening

		// 	// msg.ChannelID = msgDetail.Message.ChannelID
		// 	// msg.Message = msgDetail.Message

		// 	var msgForBroadcast wsmodels.WSChannelListItemForBroadcast
		// 	msgForBroadcast.Channel = utils.GetChannelInfo(msgDetail.Message.ChannelID)

		// 	var lastMessage db_models.Message
		// 	if err := database.DB.Where("id = ?", msgForBroadcast.Channel.LastMessageID).First(&lastMessage).Error; err != nil {

		// 	}
		// 	msgForBroadcast.LastMessage = lastMessage
		// 	msgForBroadcast.LastMessageSender = utils.GetUserInfo(1, msgForBroadcast.LastMessage.SenderID)

		// 	// query all channel member info
		// 	query := `
		// 		SELECT *
		// 		FROM channel_members
		// 		WHERE channel_id = ?
		// 		AND user_id != ?;
		// 	`

		// 	var channelMembers []db_models.ChannelMember
		// 	err := database.DB.Raw(query, msgDetail.Message.ChannelID, 1).Scan(&channelMembers).Error
		// 	if err != nil {
		// 		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channels"})
		// 		// return
		// 	}

		// 	msgForBroadcast.Users = make([]responce_models.GetUserInfoResponce, len(channelMembers))
		// 	for j := range channelMembers {
		// 		msgForBroadcast.Users[j] = utils.GetUserInfo(1, msgForBroadcast.LastMessage.SenderID)
		// 	}

		// 	for chatListClient := range chatListChannels[chatListChannelID] {
		// 		err := chatListClient.WriteJSON(msgForBroadcast)
		// 		if err != nil {
		// 			log.Printf("error: %v", err)
		// 			chatListClient.Close()
		// 			delete(chatListChannels[chatListChannelID], chatListClient)
		// 			if len(chatListChannels[chatListChannelID]) == 0 {
		// 				delete(chatListChannels, chatListChannelID)
		// 			}
		// 		}
		// 	}
		// }
	}
}
