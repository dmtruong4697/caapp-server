package routes

import (
	"caapp-server/src/controllers"
	"caapp-server/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/user-friend", middlewares.AuthMiddleware(), controllers.GetUserFriend)
	}
}
