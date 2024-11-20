package routes

import (
	wsroutes "caapp-server/src/routes/ws.routes"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Auth routes
	SetupAuthRoutes(r)

	// Profile routes
	SetupProfileRoutes(r)

	// Search routes
	SetupSearchRoutes(r)

	// Friend routes
	SetupFriendRoutes(r)

	// Chat routes
	SetupChatRoutes(r)

	// Channel routes
	SetupChannelRoutes(r)

	wsroutes.SetupWSChatListRoutes(r)

	return r
}
