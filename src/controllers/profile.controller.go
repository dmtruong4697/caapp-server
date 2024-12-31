package controllers

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	responce_models "caapp-server/src/models/responce_models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProfileInfo(c *gin.Context) {

	// var req request_models.GetProfileInfoRequest
	// if err := c.BindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request info"})
	// 	return
	// }
	userID := c.MustGet("id").(uint)

	var dbUser db_models.User
	if err := database.DB.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_022001"})
		return
	}

	var res responce_models.GetProfileInfoResponce
	res.Profile = dbUser

	c.JSON(http.StatusOK, res)
}

func UpdateProfileInfo(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.UpdateProfileRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_028001"})
		return
	}

	// check duplicate hashtag name
	var existingUser db_models.User
	if err := database.DB.Where("hashtag_name = ?", req.HashtagName).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_028002"})
	}

	// update database
	var user db_models.User
	if err := database.DB.Where("id = ?", currentUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_028003"})
		return
	}

	user.PhoneNumber = req.PhoneNumber
	user.FirstName = req.FirstName
	user.MiddleName = req.MiddleName
	user.LastName = req.LastName
	user.DateOfBirth = req.DateOfBirth
	user.HashtagName = req.HashtagName
	user.Gender = req.Gender
	user.Language = req.Language
	user.Country = req.Country
	user.ProfileDescription = req.ProfileDescription
	user.AvatarImage = req.AvatarImage
	user.CoverImage = req.CoverImage
	user.JobName = req.JobName
	user.TimeZone = req.TimeZone
	user.LastUpdate = time.Now()

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_028004"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UpdatePassword(c *gin.Context) {
	userID := c.MustGet("id").(uint)

	var dbUser db_models.User
	if err := database.DB.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var updatedPassword request_models.UpdatedPaswordRequest
	if err := c.BindJSON(&updatedPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Update account password
	dbUser.Password = updatedPassword.Password
	dbUser.LastUpdate = time.Now()

	// Save update
	if err := database.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user information"})
		return
	}

	c.JSON(http.StatusOK, dbUser)
}

func FirstUpdateProfileInfo(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.FirstUpdateProfileInfoRequest
	if err := c.BindJSON((&req)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_024001"})
		return
	}

	var user db_models.User
	if err := database.DB.Where("id = ?", currentUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_024002"})
		return
	}

	user.FirstName = req.FirstName
	user.MiddleName = req.MiddleName
	user.LastName = req.LastName
	user.PhoneNumber = req.PhoneNumber
	user.HashtagName = req.HashtagName
	user.Gender = req.Gender
	user.DateOfBirth = req.DateOfBirth
	user.Country = req.Country
	user.Language = req.Language
	user.AccountStatus = "1"

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_024003"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}

func CheckDuplicateHashtagName(c *gin.Context) {
	var req request_models.CheckDuplicateHashtagNameRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_023001"})
		return
	}

	var res responce_models.CheckDuplicateHashtagNameResponse

	var existingUser db_models.User
	if err := database.DB.Where("hashtag_name = ?", req.HashtagName).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.IsAvailableHashtagName = true
			c.JSON(http.StatusOK, res)
			return
		} else {
			// Xử lý các lỗi khác ngoài ErrRecordNotFound
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_023002"})
			return
		}
	} else {
		res.IsAvailableHashtagName = false
	}

	c.JSON(http.StatusOK, res)
}
