package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	response_models "caapp-server/src/models/responce_models"
	utils "caapp-server/src/utils"
)

func GetUserInfo(c *gin.Context) {
	currentUserID := c.MustGet("id").(int)

	var userID uint
	if err := c.BindJSON(&userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	var dbUser db_models.User
	if err := database.DB.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var friend db_models.Friend
	database.DB.Where("(first_user_id = ? AND second_user_id = ?) OR (second_user_id = ? AND first_user_id = ?)", currentUserID, userID, userID, currentUserID).First(&friend)
	var request db_models.FriendRequest
	database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", currentUserID, userID, userID, currentUserID).First(&request)

	response := response_models.GetUserInfoResponce{
		ID:                 int(dbUser.ID),
		AccountStatus:      dbUser.AccountStatus,
		AvatarImage:        dbUser.AvatarImage,
		Country:            dbUser.Country,
		CoverImage:         dbUser.CoverImage,
		CreateAt:           dbUser.CreateAt,
		DateOfBirth:        dbUser.DateOfBirth,
		Email:              dbUser.Email,
		FirstName:          dbUser.FirstName,
		MiddleName:         dbUser.MiddleName,
		LastName:           dbUser.LastName,
		Gender:             dbUser.Gender,
		HashtagName:        dbUser.HashtagName,
		JobName:            dbUser.JobName,
		Language:           dbUser.Language,
		VerificationStatus: dbUser.VerificationStatus,
		TimeZone:           dbUser.TimeZone,
		ProfileDescription: dbUser.ProfileDescription,
		PhoneNumber:        dbUser.PhoneNumber,
		LastActive:         dbUser.LastActive,
		LastUpdate:         dbUser.LastUpdate,
		Friend:             friend,
		Request:            request,
	}

	c.JSON(http.StatusOK, response)
}

func GetUserFriend(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.GetUserFriendRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_026001"})
		return
	}

	var dbUser db_models.User
	if err := database.DB.Where("id = ?", req.UserID).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_026002"})
		return
	}

	var userFriend []db_models.Friend
	if err := database.DB.Where("first_user_id = ? OR second_user_id = ?", dbUser.ID, dbUser.ID).Find(&userFriend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_026003"})
		return
	}

	var res response_models.GetUserFriendResponse
	for i := range userFriend {
		if userFriend[i].FirstUserID != dbUser.ID {
			res.Users[i] = utils.GetUserInfo(currentUserID, userFriend[i].FirstUserID)
		} else {
			res.Users[i] = utils.GetUserInfo(currentUserID, userFriend[i].SecondUserID)
		}
	}

	c.JSON(http.StatusOK, res)
}
