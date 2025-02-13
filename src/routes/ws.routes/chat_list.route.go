package wsroutes

import (
	"caapp-server/src/ws"

	"github.com/gin-gonic/gin"
)

func SetupWSChatListRoutes(r *gin.Engine) {
	wsChatListRoutes := r.Group("/ws-chatlist")
	{
		wsChatListRoutes.GET("/connect", ws.HandleChannelListConnections)
	}
}
