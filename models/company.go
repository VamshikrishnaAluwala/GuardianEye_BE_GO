package models

import (
	"time"
)

type Company struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	ApiKey    string    `gorm:"unique;not null" json:"-"` // Hide from JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
