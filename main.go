package main

import (
	"log"

	"amass-scanner/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Define routes
	router.GET("/scan", handlers.ScanHandler)
	router.POST("/batchscan", handlers.BatchScanHandler)

	// Start the server
	log.Println("Server started on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
