package routes

import (
	"caapp-server/src/controllers"
	"caapp-server/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/validate-email", controllers.ValidateEmail)
		authRoutes.POST("/login", controllers.Login)
		authRoutes.POST("/logout", middlewares.AuthMiddleware(), controllers.Logout)
	}
}
