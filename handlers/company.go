package handlers

import (
	"amass-scanner/models"
	"amass-scanner/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateCompanyRequest struct {
	Name string `json:"name" binding:"required"`
}

func CreateCompanyHandler(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company name is required"})
		return
	}

	company, err := services.CreateCompany(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company: " + err.Error()})
		return
	}

	// Return the company info including the API key
	c.JSON(http.StatusCreated, gin.H{
		"id":         company.ID,
		"name":       company.Name,
		"api_key":    company.ApiKey, // Only shown during creation
		"created_at": company.CreatedAt,
	})
}

func GetCompanyHandler(c *gin.Context) {
	// Get company from context (set by auth middleware)
	company, exists := c.Get("company")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func UpdateCompanyHandler(c *gin.Context) {
	companyID := c.MustGet("company_id").(uint)

	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company name is required"})
		return
	}

	// Update company in database
	if err := services.DB.Model(&models.Company{}).Where("id = ?", companyID).Update("name", req.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company updated successfully"})
}

func RegenerateApiKeyHandler(c *gin.Context) {
	companyID := c.MustGet("company_id").(uint)

	// Generate new API key
	newApiKey, err := services.GenerateApiKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new API key"})
		return
	}

	// Update company with new API key
	if err := services.DB.Model(&models.Company{}).Where("id = ?", companyID).Update("api_key", newApiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "API key regenerated successfully",
		"new_api_key": newApiKey,
	})
}
