package rcroutes

import (
	rcws "caapp-server/src/ws/rc_ws"

	"github.com/gin-gonic/gin"
)

func SetupRCQueueRoutes(r *gin.Engine) {
	rcQueueRoutes := r.Group("/rc")
	{
		rcQueueRoutes.GET("/queue", rcws.HandleWaitingQueueConnections)
	}
}
