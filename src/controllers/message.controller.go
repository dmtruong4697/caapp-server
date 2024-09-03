package controllers

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	responce_models "caapp-server/src/models/responce_models"
	"caapp-server/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetChannelMessage(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.GetChannelMessageRequest
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	var channelMessage responce_models.ChannelMessages
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
