package rcroutes

import (
	rcws "caapp-server/src/ws/rc_ws"

	"github.com/gin-gonic/gin"
)

func SetupRCChatRoutes(r *gin.Engine) {
	rcQueueRoutes := r.Group("/rc")
	{
		rcQueueRoutes.GET("/chat", rcws.HandleRCChatConnections)
	}
}
