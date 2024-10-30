package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	responce_models "caapp-server/src/models/responce_models"
	response_models "caapp-server/src/models/responce_models"

	utils "caapp-server/src/utils"
)

func GetChannelList(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var getChannelListRequest request_models.GetChannelListRequest
	if err := c.BindJSON(&getChannelListRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request"})
		return
	}

	query := `
		SELECT channels.*
		FROM channels
		JOIN channel_members ON channels.id = channel_members.channel_id
		WHERE channel_members.user_id = ?
		LIMIT ?
		OFFSET ?;
	`
	var channels []db_models.Channel
	err := database.DB.Select(&channels, query, currentUserID, getChannelListRequest.Limit, getChannelListRequest.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channels"})
		return
	}

	var res response_models.GetChannelListResponce
	for i := range channels {
		var lastMessage db_models.Message
		if err := database.DB.Where("id = ?", channels[i].LastMessageID).First(&lastMessage).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
			return
		}

		res.Channels[i].Channel = channels[i]
		res.Channels[i].LastMessage = lastMessage
		res.Channels[i].LastMessageSender = utils.GetUserInfo(currentUserID, lastMessage.SenderID)
	}
	res.IsLastPage = len(channels) < getChannelListRequest.Limit

	c.JSON(http.StatusOK, res)
}

func CheckFriendchannel(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var checkFriendChannelRequest request_models.CheckFriendChannelRequest
	if err := c.BindJSON(&checkFriendChannelRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request"})
		return
	}

	query := `
		SELECT channels.*
		FROM channels
		JOIN channel_members cm1 ON channels.id = cm1.channel_id
		JOIN channel_members cm2 ON channels.id = cm2.channel_id
		WHERE channels.type = 'friend'
		AND cm1.user_id = ?
		AND cm2.user_id = ?
		AND cm1.channel_id = cm2.channel_id
	`

	var channel db_models.Channel
	result := database.DB.Raw(query, currentUserID, checkFriendChannelRequest.UserID).First(&channel)
	if result.Error != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channel"})
		// return
	}

	if channel.ID == 0 {
		newChannel := db_models.Channel{
			CreatorID: currentUserID,
			CreateAt:  time.Now(),
			Type:      "friend",
		}

		if err := database.DB.Create(&newChannel).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel"})
			return
		}

		channelMember1 := db_models.ChannelMember{
			UserID:    currentUserID,
			ChannelID: newChannel.ID,
			JoinAt:    time.Now(),
			Role:      "member",
		}

		channelMember2 := db_models.ChannelMember{
			UserID:    checkFriendChannelRequest.UserID,
			ChannelID: newChannel.ID,
			JoinAt:    time.Now(),
			Role:      "member",
		}

		if err := database.DB.Create(&channelMember1).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add current user to channel"})
			return
		}
		if err := database.DB.Create(&channelMember2).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add friend user to channel"})
			return
		}

		var res response_models.CheckFriendChannelResponce
		res.Channel = newChannel

		c.JSON(http.StatusOK, res)
		return
	}

	var res response_models.CheckFriendChannelResponce
	res.Channel = channel

	c.JSON(http.StatusOK, res)
}

func GetFriendChannelInfo(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var getFriendChannelInfoRequest request_models.GetFriendChannelInfoRequest
	if err := c.BindJSON(&getFriendChannelInfoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request"})
		return
	}

	var channel db_models.Channel
	if err := database.DB.Where("id = ?", getFriendChannelInfoRequest.ChannelID).First(&channel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	query := `
		SELECT *
		FROM channel_members
		WHERE channel_id = ?
		AND user_id != ?;
	`

	var channelMember db_models.ChannelMember
	result := database.DB.Raw(query, getFriendChannelInfoRequest.ChannelID, currentUserID).First(&channelMember)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channel member"})
		return
	}

	var res response_models.GetFriendChannelInfoResponce
	res.Channel = channel
	res.User = utils.GetUserInfo(currentUserID, channelMember.UserID)

	c.JSON(http.StatusOK, res)
}

func GetGroupChannelInfo(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var getGroupChannelInfoRequest request_models.GetGroupChannelInfoRequest
	if err := c.BindJSON(&getGroupChannelInfoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request"})
		return
	}

	var channel db_models.Channel
	if err := database.DB.Where("id = ?", getGroupChannelInfoRequest.ChannelID).First(&channel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
		return
	}

	query := `
		SELECT *
		FROM channel_members
		WHERE channel_id = ?
		AND user_id != ?;
	`

	var channelMembers []db_models.ChannelMember
	err := database.DB.Select(&channelMembers, query, getGroupChannelInfoRequest.ChannelID, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channel member"})
		return
	}

	var res response_models.GetGroupChannelInfoResponce
	res.Channel = channel
	for i := range channelMembers {
		res.Users[i] = utils.GetUserInfo(currentUserID, channelMembers[i].UserID)
	}

	c.JSON(http.StatusOK, res)
}

func GetChannelChatHistory(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.GetChannelChatHistoryRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request info"})
		return
	}

	var channelMember db_models.ChannelMember
	if err := database.DB.Where("channel_id = ? AND user_id = ?", req.ChannelID, currentUserID).First(&channelMember).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not a member of the channel"})
		return
	}

	var messages []db_models.Message
	if err := database.DB.Where("channel_id = ?", req.ChannelID).Find(&messages).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		// return
	}

	var channelMessage responce_models.GetChannelChatHistoryResponce
	for i := range messages {
		var medias []db_models.Media
		if err := database.DB.Where("message_id = ?", messages[i].ID).Find(&medias).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch medias"})
			return
		}

		channelMessage.Messages[i].Message = messages[i]
		channelMessage.Messages[i].Sender = utils.GetUserInfo(currentUserID, messages[i].SenderID)
		channelMessage.Messages[i].Media = medias
	}

	c.JSON(http.StatusOK, channelMessage)
}
