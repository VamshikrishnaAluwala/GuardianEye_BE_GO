package utils

import (
	"bufio"
	"mime/multipart"
	"strings"
)

// ReadDomainsFromFile reads domains from an uploaded file
func ReadDomainsFromFile(file multipart.File) ([]string, error) {
	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			domains = append(domains, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return domains, nil
}
