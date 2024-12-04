package rccontrollers

import (
	"caapp-server/src/database"
	models "caapp-server/src/models/db_models"
	rcdbmodels "caapp-server/src/models/db_models/rc_db_models"
	rcrequestmodels "caapp-server/src/models/request_models/rc_request_models"
	rcresponsemodel "caapp-server/src/models/responce_models/rc_response_model"
	"caapp-server/src/utils/helper"
	rcws "caapp-server/src/ws/rc_ws"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRCChannelChatHistory(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req rcrequestmodels.GetRCChannelChatHistoryRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request info"})
		return
	}

	var channelMember rcdbmodels.RCChannelMember
	if err := database.DB.Where("channel_id = ? AND user_id = ?", req.ChannelID, currentUserID).First(&channelMember).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not a member of the channel"})
		return
	}

	var messages []rcdbmodels.RCMessage
	if err := database.DB.Where("channel_id = ?", req.ChannelID).Order("create_at DESC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	var channelMessage rcresponsemodel.GetRCChannelChatHistoryResponce

	channelMessage.Messages = make([]rcresponsemodel.GetRCChannelChatHistoryItem, len(messages))

	for i := range messages {
		var medias []rcdbmodels.RCMedia
		if err := database.DB.Where("message_id = ?", messages[i].ID).Find(&medias).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch medias"})
			return
		}

		channelMessage.Messages[i].Message = messages[i]
		channelMessage.Messages[i].Media = medias
	}

	c.JSON(http.StatusOK, channelMessage)
}

func GetCurrentRCChannel(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var dbUser models.User
	if err := database.DB.Where("id = ?", currentUserID).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "khong tim thay nguoi dung"})
		return
	}

	var RCChannel rcdbmodels.RCChannel
	if err := database.DB.Where("id = ?", dbUser.CurrentRCChannelID).First(&RCChannel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Khong tim thay channel"})
		return
	}

	c.JSON(http.StatusOK, RCChannel)
}

func LeaveCurrentRCChannel(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var dbUser models.User
	if err := database.DB.Where("id = ?", currentUserID).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "khong tim thay nguoi dung"})
		return
	}

	var RCChannel rcdbmodels.RCChannel
	if err := database.DB.Where("id = ?", dbUser.CurrentRCChannelID).First(&RCChannel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Khong tim thay channel"})
		return
	}

	var channelMember rcdbmodels.RCChannelMember
	if err := database.DB.Where("user_id = ? AND channel_id = ?", dbUser.ID, RCChannel.ID).First(&channelMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Khong tim thay channel member"})
		return
	}

	var newNotificationMessage rcdbmodels.RCMessage
	newNotificationMessage.ChannelID = RCChannel.ID
	newNotificationMessage.CreateAt = time.Now()
	newNotificationMessage.LastUpdate = time.Now()
	newNotificationMessage.Type = "2"

	// create new channel notification message
	if err := database.DB.Create(&newNotificationMessage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "loi tao moi tin nhan"})
		return
	}

	// update channel last message id
	RCChannel.LastMessageID = newNotificationMessage.ID
	if err := database.DB.Save(&RCChannel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "loi update channel"})
		return
	}

	// update user's current rc channel
	dbUser.CurrentRCChannelID = 0
	if err := database.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "loi update db user"})
		return
	}

	// delete rc channel member
	if err := database.DB.Delete(&channelMember).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "loi xoa channel member"})
		return
	}

	// send new notification message to ws chat channel
	chatChannelID := helper.UIntToString(RCChannel.ID) + "-RCChatChannel"
	for client := range rcws.GetRCChatChannels()[chatChannelID] {
		err := client.WriteJSON(newNotificationMessage)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(rcws.GetRCChatChannels()[chatChannelID], client)
			if len(rcws.GetRCChatChannels()[chatChannelID]) == 0 {
				delete(rcws.GetRCChatChannels(), chatChannelID)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}
