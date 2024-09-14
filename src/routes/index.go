package routes

import (
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

	return r
}
