package rcroutes

import (
	rccontrollers "caapp-server/src/controllers/rc_controllers"
	"caapp-server/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRCChannelRoutes(r *gin.Engine) {
	rcChannelRoutes := r.Group("/rc-channel")
	{
		rcChannelRoutes.POST("/chat-history", middlewares.AuthMiddleware(), rccontrollers.GetRCChannelChatHistory)
		rcChannelRoutes.POST("/current-rc-channel", middlewares.AuthMiddleware(), rccontrollers.GetCurrentRCChannel)
	}
}
