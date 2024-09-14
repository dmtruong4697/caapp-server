package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"caapp-server/src/database"
	db_models "caapp-server/src/models/db_models"
	responce_models "caapp-server/src/models/responce_models"
	utils "caapp-server/src/utils"
)

func SearchUserByHashtagName(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)
	// hastagName := c.Query("hashtag_name")

	var users []db_models.User
	database.DB.Find(&users)

	var res responce_models.SearchUserByHashtagNameResponce
	for i := range users {
		res.Users[i] = utils.GetUserInfo(currentUserID, users[i].ID)
	}

	c.JSON(http.StatusOK, res)
}
