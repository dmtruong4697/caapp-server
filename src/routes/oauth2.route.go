package routes

import (
	"caapp-server/src/controllers/oauth2"

	"github.com/gin-gonic/gin"
)

func SetupOAuth2Routes(r *gin.Engine) {
	authRoutes := r.Group("/oauth2")
	{
		authRoutes.POST("/google-login", oauth2.HandleGoogleLogin)
	}
}
