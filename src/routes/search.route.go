package routes

import (
	"caapp-server/src/controllers"
	"caapp-server/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupSearchRoutes(r *gin.Engine) {
	profileRoutes := r.Group("/search")
	{
		profileRoutes.POST("/hashtag-name", middlewares.AuthMiddleware(), controllers.SearchUserByHashtagName)
	}
}
