package routes

import (
	"caapp-server/src/controllers"
	"caapp-server/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupProfileRoutes(r *gin.Engine) {
	profileRoutes := r.Group("/profile")
	{
		profileRoutes.POST("/profile-info", middlewares.AuthMiddleware(), controllers.GetProfileInfo)
		profileRoutes.POST("/check-duplicate-hashtag-name", controllers.CheckDuplicateHashtagName)
		profileRoutes.POST("/first-update-profile-info", middlewares.AuthMiddleware(), controllers.FirstUpdateProfileInfo)
		profileRoutes.POST("/update-profile-info", middlewares.AuthMiddleware(), controllers.UpdateProfileInfo)
	}
}
