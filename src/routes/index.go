package routes

import (
	rcroutes "caapp-server/src/routes/rc.routes"
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

	// Constant Data routes
	SetupConstantDataRoutes(r)

	wsroutes.SetupWSChatListRoutes(r)

	// rc module
	rcroutes.SetupRCQueueRoutes(r)
	rcroutes.SetupRCChannelRoutes(r)
	rcroutes.SetupRCChatRoutes(r)

	return r
}
