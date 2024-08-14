package controllers

import (
	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProfileInfo(c *gin.Context) {
	userID := c.MustGet("id").(uint)

	var dbUser db_models.User
	if err := database.DB.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, dbUser)
}

func UpdateProfileInfo(c *gin.Context) {
	email := c.MustGet("email").(string)

	var dbUser db_models.User
	if err := database.DB.Where("email = ?", email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var updatedProfile request_models.UpdatedProfileRequest
	if err := c.BindJSON(&updatedProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Update user profile
	dbUser.FirstName = updatedProfile.FirstName
	dbUser.MiddleName = updatedProfile.MiddleName
	dbUser.LastName = updatedProfile.LastName
	dbUser.PhoneNumber = updatedProfile.PhoneNumber
	dbUser.AvatarImage = updatedProfile.AvatarImage
	dbUser.CoverImage = updatedProfile.CoverImage
	dbUser.HashtagName = updatedProfile.HashtagName
	dbUser.Gender = updatedProfile.Gender
	dbUser.DateOfBirth = updatedProfile.DateOfBirth
	dbUser.Language = updatedProfile.Language
	dbUser.Country = updatedProfile.Country
	dbUser.ProfileDescription = updatedProfile.ProfileDescription
	dbUser.JobName = updatedProfile.JobName
	dbUser.TimeZone = updatedProfile.TimeZone
	dbUser.LastUpdate = time.Now()

	// Save update
	if err := database.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user information"})
		return
	}

	c.JSON(http.StatusOK, dbUser)
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
