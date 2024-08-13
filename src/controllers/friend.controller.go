package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"caapp-server/src/database"
	"caapp-server/src/enums"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	responce_models "caapp-server/src/models/responce_models"
)

func GetRelationship(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.GetRelationshipRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	var friend db_models.Friend
	database.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		currentUserID,
		req.UserID,
		req.UserID,
		currentUserID,
	).Find(&friend)

	res := responce_models.GetRelationshipResponse{Relationship: friend}
	if friend.ID > 0 {
		res.Status = enums.RELATIONSHIP_FRIEND
	} else {
		res.Status = enums.RELATIONSHIP_NOT_FRIEND
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
