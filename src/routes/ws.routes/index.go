package wsroutes

import (
	"github.com/gin-gonic/gin"
)

func SetupWSRouter() *gin.Engine {
	r := gin.Default()

	// Chat list routes
	SetupWSChatListRoutes(r)

	return r
}
