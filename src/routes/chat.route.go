package routes

import (
	"caapp-server/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(r *gin.Engine) {
	chatRoutes := r.Group("/chat")
	{
		chatRoutes.GET("/ws", controllers.HandleConnections)
	}
}
