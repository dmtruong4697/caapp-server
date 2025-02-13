package controllers

import (
	"net/http"

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
		"(first_user_id_id = ? AND second_user_id = ?) OR (second_user_id = ? AND first_user_id_id = ?)",
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
