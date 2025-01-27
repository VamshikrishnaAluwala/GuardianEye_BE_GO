package services

import (
	"amass-scanner/models"
	"crypto/rand"
	"encoding/hex"
)

func GenerateApiKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func CreateCompany(name string) (*models.Company, error) {
	apiKey, err := GenerateApiKey()
	if err != nil {
		return nil, err
	}

	company := &models.Company{
		Name:   name,
		ApiKey: apiKey,
	}

	if err := DB.Create(company).Error; err != nil {
		return nil, err
	}

	return company, nil
}

func GetCompanyByApiKey(apiKey string) (*models.Company, error) {
	var company models.Company
	if err := DB.Where("api_key = ?", apiKey).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}
