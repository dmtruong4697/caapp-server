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
		profileRoutes.POST("/request", middlewares.AuthMiddleware(), controllers.GetFriendRequest)
		profileRoutes.POST("/received-request", middlewares.AuthMiddleware(), controllers.GetAllFriendRequestReceived)
		profileRoutes.POST("/sent-request", middlewares.AuthMiddleware(), controllers.GetAllFriendRequestSent)
		profileRoutes.POST("/create-request", middlewares.AuthMiddleware(), controllers.CreateFriendRequest)
		profileRoutes.POST("/accept-request", middlewares.AuthMiddleware(), controllers.AcceptFriendRequest)
		profileRoutes.POST("/reject-request", middlewares.AuthMiddleware(), controllers.RejectFriendRequest)
		profileRoutes.POST("/delete-request", middlewares.AuthMiddleware(), controllers.DeleteFriendRequest)

		profileRoutes.POST("/all-my-friend", middlewares.AuthMiddleware(), controllers.GetAllMyFriend)
	}
}
