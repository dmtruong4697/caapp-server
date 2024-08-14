package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	responce_models "caapp-server/src/models/responce_models"

	utils "caapp-server/src/utils"
)

func GetFriendRequest(c *gin.Context) {
	var req request_models.GetFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request id info"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
		return
	}

	res := responce_models.GetFriendRequestResponce{FriendRequest: friendRequest}

	c.JSON(http.StatusOK, res)
}

func GetAllFriendRequestReceived(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var requests []db_models.FriendRequest
	database.DB.Where("receiver_id = ?", currentUserID).Find(&requests)

	var res responce_models.GetListFriendRequestReceivedResponce
	for i := range requests {
		res.Requests[i].User = utils.GetUserInfo(currentUserID, requests[i].SenderID)
		res.Requests[i].FriendRequest = requests[i]
	}

	c.JSON(http.StatusOK, res)
}

func GetAllFriendRequestSent(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var requests []db_models.FriendRequest
	database.DB.Where("sender_id = ?", currentUserID).Find(&requests)

	var res responce_models.GetListFriendRequestSentResponce
	for i := range requests {
		res.Requests[i].User = utils.GetUserInfo(currentUserID, requests[i].ReceiverID)
		res.Requests[i].FriendRequest = requests[i]
	}

	c.JSON(http.StatusOK, res)
}

func CreateFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.CreateFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	friend_request := db_models.FriendRequest{
		SenderID:   currentUserID,
		ReceiverID: req.UserID,
		CreateAt:   time.Now(),
	}

	if err := database.DB.Save(&friend_request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request created successfully"})
}

func AcceptFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.AcceptFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request info"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
		return
	}

	friend := db_models.Friend{
		FirstUserID:  currentUserID,
		SecondUserID: friendRequest.SenderID,
		CreateAt:     time.Now(),
	}

	if err := database.DB.Create(&friend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create friend"})
		return
	}

	if err := database.DB.Delete(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}

func RejectFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.RefuseFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request info"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
		return
	}

	if friendRequest.ReceiverID != currentUserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to reject this friend request"})
		return
	}

	if err := database.DB.Delete(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request rejected successfully"})
}

func DeleteFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.DeleteFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request info"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
		return
	}

	if friendRequest.SenderID != currentUserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete this friend request"})
		return
	}

	if err := database.DB.Delete(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request delete successfully"})
}
