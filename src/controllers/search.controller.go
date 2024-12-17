package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	request_models "caapp-server/src/models/request_models"
	responce_models "caapp-server/src/models/responce_models"
	utils "caapp-server/src/utils"
)

func SearchUserByHashtagName(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req request_models.SearchUserByHashtagNameRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_025001"})
		return
	}

	var res responce_models.SearchUserByHashtagNameResponse

	var dbUsers db_models.User
	if err := database.DB.Where("hashtag_name = ?", req.HashtagName).First(&dbUsers).Error; err != nil {
		res.IsFound = false
		c.JSON(http.StatusOK, res)
	}

	res.User = utils.GetUserInfo(currentUserID, dbUsers.ID)
	res.IsFound = true

	c.JSON(http.StatusOK, res)
}
