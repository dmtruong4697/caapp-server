package routes

import (
	"caapp-server/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupConstantDataRoutes(r *gin.Engine) {
	authRoutes := r.Group("/constant-data")
	{
		authRoutes.POST("/languages", controllers.GetLanguageDataList)
	}
}
