package services

import (
	"os/exec"
	"strings"
)

// RunAmassScan runs the Amass scan for a given domain
func RunAmassScan(domain string) ([]string, error) {
	cmd := exec.Command("amass", "enum", "-d", domain)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseAmassOutput(string(output)), nil
}

// parseAmassOutput processes the Amass output into a slice of results
func parseAmassOutput(output string) []string {
	lines := strings.Split(output, "\n")
	var results []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			results = append(results, strings.TrimSpace(line))
		}
	}
	return results
}
