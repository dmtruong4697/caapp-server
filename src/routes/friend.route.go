package routes

import (
	"caapp-server/src/controllers"
	"caapp-server/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupFriendRoutes(r *gin.Engine) {
	profileRoutes := r.Group("/friend")
	{
		profileRoutes.POST("/suggest", middlewares.AuthMiddleware(), controllers.GetSuggestUser)
	}
}
