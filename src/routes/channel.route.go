package routes

import (
	"caapp-server/src/controllers"
	"caapp-server/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupChannelRoutes(r *gin.Engine) {
	channelRoutes := r.Group("/channel")
	{
		channelRoutes.POST("/channel-list", middlewares.AuthMiddleware(), controllers.GetChannelList)
		channelRoutes.POST("/check-friend-channel", middlewares.AuthMiddleware(), controllers.CheckFriendchannel)
		channelRoutes.POST("/friend-channel-info", middlewares.AuthMiddleware(), controllers.GetFriendChannelInfo)
		channelRoutes.POST("/group-channel-info", middlewares.AuthMiddleware(), controllers.GetGroupChannelInfo)
	}
}
