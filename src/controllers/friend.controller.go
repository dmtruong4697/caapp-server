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
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000017"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error_code": "api_error_404_000018"})
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
	res.Requests = make([]responce_models.GetListFriendRequestReceivedResponceItem, len(requests))

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
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000019"})
		return
	}

	friend_request := db_models.FriendRequest{
		SenderID:   currentUserID,
		ReceiverID: req.UserID,
		CreateAt:   time.Now(),
	}

	if err := database.DB.Save(&friend_request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000020"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request created successfully"})
}

func AcceptFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.AcceptFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000021"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error_code": "api_error_404_000022"})
		return
	}

	friend := db_models.Friend{
		FirstUserID:  currentUserID,
		SecondUserID: friendRequest.SenderID,
		CreateAt:     time.Now(),
	}

	if err := database.DB.Create(&friend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000023"})
		return
	}

	if err := database.DB.Delete(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000024"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}

func RejectFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.RefuseFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000025"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "api_error_404_000026"})
		return
	}

	if friendRequest.ReceiverID != currentUserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "api_error_401_000027"})
		return
	}

	if err := database.DB.Delete(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000028"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request rejected successfully"})
}

func DeleteFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.DeleteFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_000029"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error_code": "api_error_404_000030"})
		return
	}

	if friendRequest.SenderID != currentUserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "api_error_401_000031"})
		return
	}

	if err := database.DB.Delete(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_000032"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request delete successfully"})
}

func GetSuggestUser(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var users []db_models.User
	database.DB.Find(&users)

	var res responce_models.SearchUserByHashtagNameResponce
	res.Users = make([]responce_models.GetUserInfoResponce, len(users))

	for i := range users {
		res.Users[i] = utils.GetUserInfo(currentUserID, users[i].ID)
	}

	c.JSON(http.StatusOK, res)
}

func GetAllMyFriend(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var friends []db_models.Friend
	database.DB.Where(
		"(first_user_id = ?) OR (second_user_id = ?)",
		currentUserID,
		currentUserID,
	).Find(&friends)

	var res responce_models.GetAllMyFriendResponce
	res.Friends = make([]responce_models.GetUserInfoResponce, len(friends))

	for i := range friends {
		var friendId uint
		if currentUserID == friends[i].FirstUserID {
			friendId = friends[i].SecondUserID
		} else {
			friendId = friends[i].FirstUserID
		}
		res.Friends[i] = utils.GetUserInfo(currentUserID, friendId)
	}

	c.JSON(http.StatusOK, res)
}

func GetAllUserFriend(c *gin.Context) {

}
