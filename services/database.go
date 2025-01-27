package services

import (
	"amass-scanner/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=Vamshi@304 dbname=guardianeye_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.ScanResult{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = db
}
