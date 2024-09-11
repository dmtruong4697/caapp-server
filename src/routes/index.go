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

	return r
}
