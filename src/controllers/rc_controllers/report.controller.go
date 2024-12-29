package rccontrollers

import (
	"caapp-server/src/database"
	rcdbmodels "caapp-server/src/models/db_models/rc_db_models"
	rcrequestmodels "caapp-server/src/models/request_models/rc_request_models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ReportRCChannel(c *gin.Context) {
	currentUserID := c.MustGet("id").(uint)

	var req rcrequestmodels.RCReportRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "api_error_400_027001"})
		return
	}

	var newReport rcdbmodels.RCReport
	newReport.SenderID = currentUserID
	newReport.ChannelID = req.ChannelID
	newReport.Content = req.Content
	newReport.Type = req.Type
	newReport.CreateAt = time.Now()

	if err := database.DB.Create(&newReport).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "api_error_500_027002"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
