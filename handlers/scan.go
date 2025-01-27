package handlers

import (
	"amass-scanner/models"
	"amass-scanner/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ScanHandler(c *gin.Context) {
	companyID := c.MustGet("company_id").(uint)

	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'domain' query parameter"})
		return
	}

	var params models.ScanParams
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters format"})
		return
	}

	if len(params.CommandArgs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "command_args cannot be empty"})
		return
	}

	// Run the scan
	result, err := services.RunAmassScan(domain, params, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch subdomains associated with the scan result
	var subdomains []models.Subdomain
	if err := services.DB.Where("scan_result_id = ?", result.ID).Find(&subdomains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching subdomains"})
		return
	}

	// Extract subdomains from the database result
	subdomainResults := make([]string, len(subdomains))
	for i, subdomain := range subdomains {
		subdomainResults[i] = subdomain.Subdomain
	}

	c.JSON(http.StatusOK, gin.H{
		"domain":     domain,
		"scan_id":    result.ID,
		"subdomains": subdomainResults,
		"created_at": result.CreatedAt,
	})
}
