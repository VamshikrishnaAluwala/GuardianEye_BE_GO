package handlers

import (
	"amass-scanner/models"
	"amass-scanner/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHistoricalResultsHandler(c *gin.Context) {
	var filter models.ScanResultFilter
	if err := c.BindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters"})
		return
	}

	results, total, err := services.GetScanResults(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"total":   total,
		"limit":   filter.Limit,
		"offset":  filter.Offset,
	})
}

func GetScanResultHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing scan ID"})
		return
	}

	var scanID uint
	if _, err := fmt.Sscanf(id, "%d", &scanID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID format"})
		return
	}

	result, err := services.GetScanResultByID(scanID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scan result not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}
