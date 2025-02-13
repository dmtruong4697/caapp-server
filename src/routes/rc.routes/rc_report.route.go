package rcroutes

import (
	rccontrollers "caapp-server/src/controllers/rc_controllers"
	"caapp-server/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRCReportRoutes(r *gin.Engine) {
	rcReportRoutes := r.Group("/rc")
	{
		rcReportRoutes.GET("/report", middlewares.AuthMiddleware(), rccontrollers.ReportRCChannel)
	}
}
