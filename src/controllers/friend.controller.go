package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	responce_models "caapp-server/src/models/responce_models"

	"math/rand"

	utils "caapp-server/src/utils"
)

func GetFriendRequest(c *gin.Context) {
	var req request_models.GetFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_017001"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_017002"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_018001"})
		return
	}

	var existingFriendRequest db_models.FriendRequest
	if err := database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (receiver_id = ? AND sender_id = ?)", currentUserID, req.UserID, currentUserID, req.UserID).First(&existingFriendRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			friend_request := db_models.FriendRequest{
				SenderID:   currentUserID,
				ReceiverID: req.UserID,
				CreateAt:   time.Now(),
			}

			if err := database.DB.Save(&friend_request).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_018002"})
				return
			}

			c.JSON(http.StatusOK, gin.H{})
		} else {
			// Xử lý các lỗi khác ngoài ErrRecordNotFound
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_018003"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_500_018004"})
		return
	}
}

func AcceptFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.AcceptFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_019001"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_019002"})
		return
	}

	friend := db_models.Friend{
		FirstUserID:  currentUserID,
		SecondUserID: friendRequest.SenderID,
		CreateAt:     time.Now(),
	}

	if err := database.DB.Create(&friend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_019003"})
		return
	}

	if err := database.DB.Delete(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_019004"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}

func RejectFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.RefuseFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_020001"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "api_error_500_020002"})
		return
	}

	if friendRequest.ReceiverID != currentUserID {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_020003"})
		return
	}

	if err := database.DB.Delete(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_020004"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request rejected successfully"})
}

func DeleteFriendRequest(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.DeleteFriendRequestRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_021001"})
		return
	}

	var friendRequest db_models.FriendRequest
	if err := database.DB.First(&friendRequest, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_021002"})
		return
	}

	if friendRequest.SenderID != currentUserID {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_021003"})
		return
	}

	if err := database.DB.Delete(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_021004"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request delete successfully"})
}

type GetSuggestUserResponce struct {
	Users []responce_models.GetUserInfoResponce `json:"users"`
}

func GetSuggestUser(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var users []db_models.User
	database.DB.Find(&users)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })

	var res GetSuggestUserResponce
	res.Users = make([]responce_models.GetUserInfoResponce, 0)

	for i := range users {
		if len(res.Users) >= 5 {
			break
		}
		user := utils.GetUserInfo(currentUserID, users[i].ID)
		if user.Request.ID <= 0 {
			res.Users = append(res.Users, user)
		}
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
