package services

import (
	"amass-scanner/models"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// buildAmassCommand Creates the commanline to run
func buildAmassCommand(domain string, params models.ScanParams) *exec.Cmd {
	// Start with the base command and enum subcommand
	args := []string{"enum"}

	// Add all the command arguments from the request
	args = append(args, params.CommandArgs...)

	// Add the domain last if it's not already included
	domainIncluded := false
	for i, arg := range args {
		if arg == "-d" && i+1 < len(args) {
			domainIncluded = true
			break
		}
	}

	if !domainIncluded {
		args = append(args, "-d", domain)
	}

	log.Printf("Executing Amass command: amass %s", strings.Join(args, " "))
	return exec.Command("amass", args...)
}

// RunAmassScan runs the Amass scan for a given domain
func RunAmassScan(domain string, params models.ScanParams, companyID uint) (*models.ScanResult, error) {
	log.Printf("Starting scan for domain: %s, companyID: %d", domain, companyID)

	// Check if Amass is installed
	_, err := exec.LookPath("amass")
	if err != nil {
		log.Printf("Error: amass is not installed or not in PATH: %v", err)
		return nil, fmt.Errorf("amass is not installed or not in PATH: %v", err)
	}

	cmd := buildAmassCommand(domain, params)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	paramsJSON, _ := json.Marshal(params)
	log.Printf("Scan parameters: %s", string(paramsJSON))

	// Create the ScanResult object
	scanResult := &models.ScanResult{
		CompanyID: companyID,
		Domain:    domain,
		CreatedAt: time.Now(),
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Amass command failed: %v\nOutput: %s", err, outputStr)
		log.Printf("Error during scan: %s", errorMsg)
		scanResult.Error = errorMsg

		// Save failed scan result
		if dbErr := DB.Create(scanResult).Error; dbErr != nil {
			log.Printf("Error saving failed scan result: %v", dbErr)
		}
		return scanResult, fmt.Errorf(errorMsg)
	}

	// Parse results into subdomains
	subdomains := parseAmassOutput(outputStr)
	log.Printf("Scan completed. Found %d subdomains", len(subdomains))

	// Save the ScanResult
	if err := DB.Create(scanResult).Error; err != nil {
		log.Printf("Error saving scan result to database: %v", err)
		return scanResult, fmt.Errorf("error saving scan result: %v", err)
	}

	// Save Subdomains
	for _, subdomain := range subdomains {
		subdomainRecord := &models.Subdomain{
			ScanResultID: scanResult.ID,
			Subdomain:    subdomain,
			CreatedAt:    time.Now(),
		}
		if err := DB.Create(subdomainRecord).Error; err != nil {
			log.Printf("Error saving subdomain: %v", err)
		}
	}

	return scanResult, nil
}

// RunBatchScan Runs scan for multiple domains
func RunBatchScan(request *models.BatchScanRequest, companyID uint) (map[string]*models.ScanResult, error) {
	results := make(map[string]*models.ScanResult)

	for _, domain := range request.Domains {
		result, err := RunAmassScan(domain, request.Params, companyID)
		if err != nil {
			results[domain] = &models.ScanResult{
				CompanyID: companyID,
				Domain:    domain,
				Error:     err.Error(),
			}
		} else {
			results[domain] = result
		}
	}

	return results, nil
}

func parseAmassOutput(output string) []string {
	lines := strings.Split(output, "\n")
	var results []string
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			results = append(results, trimmed)
		}
	}
	log.Printf("Parsed %d lines from Amass output", len(results))
	return results
}

// GetScanResults gets the scan results from the database for a company
func GetScanResults(filter models.ScanResultFilter) ([]models.ScanResult, int64, error) {
	var results []models.ScanResult
	var total int64

	query := DB.Model(&models.ScanResult{}).Where("company_id = ?", filter.CompanyID)

	if filter.Domain != "" {
		query = query.Where("domain LIKE ?", "%"+filter.Domain+"%")
	}

	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate)
	}

	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetScanResultByID fetches the result for an id
func GetScanResultByID(id uint) (*models.ScanResult, error) {
	var result models.ScanResult
	if err := DB.First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
