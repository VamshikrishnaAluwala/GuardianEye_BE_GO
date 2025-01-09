package handlers

import (
	"net/http"

	"amass-scanner/services"
	"github.com/gin-gonic/gin"
)

// ScanHandler handles single domain scans
func ScanHandler(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'domain' query parameter"})
		return
	}

	results, err := services.RunAmassScan(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"domain": domain, "results": results})
}
