package handlers

import (
	"amass-scanner/models"
	"amass-scanner/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BatchScanHandler handles multiple domain scans from a file
func BatchScanHandler(c *gin.Context) {
	// Get company ID from the context (set by auth middleware)
	companyID := c.MustGet("company_id").(uint)

	var request models.BatchScanRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(request.Domains) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No domains provided"})
		return
	}

	// Pass both the request and companyID to RunBatchScan
	results, err := services.RunBatchScan(&request, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
