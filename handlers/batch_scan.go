package handlers

import (
	"net/http"

	"amass-scanner/services"
	"amass-scanner/utils"
	"github.com/gin-gonic/gin"
)

// BatchScanHandler handles multiple domain scans from a file
func BatchScanHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'file' upload"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file: " + err.Error()})
		return
	}
	defer src.Close()

	domains, err := utils.ReadDomainsFromFile(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading file: " + err.Error()})
		return
	}

	results := make(map[string][]string)
	for _, domain := range domains {
		scanResults, err := services.RunAmassScan(domain)
		if err != nil {
			results[domain] = []string{"Error: " + err.Error()}
		} else {
			results[domain] = scanResults
		}
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
