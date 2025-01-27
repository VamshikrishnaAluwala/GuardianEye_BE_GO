package main

import (
	"amass-scanner/handlers"
	"amass-scanner/middleware"
	"amass-scanner/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	services.InitDB()
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Public endpoints
	router.POST("/companies", handlers.CreateCompanyHandler)
	router.GET("/test", handlers.TestAmassHandler)

	// Protected endpoints requiring API key
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Company management
		protected.GET("/company", handlers.GetCompanyHandler)
		protected.PUT("/company", handlers.UpdateCompanyHandler)
		protected.POST("/company/regenerate-key", handlers.RegenerateApiKeyHandler)

		// Scanning endpoints
		protected.POST("/scan", handlers.ScanHandler)
		protected.POST("/batchscan", handlers.BatchScanHandler)
		protected.GET("/results", handlers.GetHistoricalResultsHandler)
		protected.GET("/results/:id", handlers.GetScanResultHandler)
		protected.GET("/nmap", handlers.NmapHandler)
		protected.GET("/livenmap", handlers.Live_NmapHandler)
	}

	log.Println("Server started on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
