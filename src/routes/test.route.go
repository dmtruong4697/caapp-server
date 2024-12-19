package routes

import (
	"caapp-server/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupTestRoutes(r *gin.Engine) {
	profileRoutes := r.Group("/test")
	{
		profileRoutes.POST("/send-mail", controllers.TestSendMail)
	}
}
