package handlers

import (
    "bufio"
    "github.com/gin-gonic/gin"
    "io"
    "log"
    "net/http"
    "os/exec"
    "strings"
)

// Live_NmapHandler handles the /nmap endpoint for multiple domains.
func Live_NmapHandler(c *gin.Context) {
    file, _, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"detail": "File parameter is required"})
        return
    }
    defer file.Close()

    domains, err := readliveDomains(file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
        return
    }

    results := make(map[string]string)
    for _, domain := range domains {
        output, err := runLiveHostNmap(domain)
        if err != nil {
            results[domain] = err.Error()
        } else {
            results[domain] = output
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "results": results,
    })
}

// readDomains reads domains from the provided file.
func readliveDomains(file io.Reader) ([]string, error) {
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

// runLiveHostNmap executes the Docker-based Nmap command for live host enumeration.
func runLiveHostNmap(domain string) (string, error) {
    // Prepare the Docker-based Nmap command
    cmd := exec.Command("docker", "run", "--rm", "instrumentisto/nmap",
        "-T5", "-sn", "-PE", "-vv", "-oX", "-", domain)
		// command = [
        //     "docker", "run", "--rm", "instrumentisto/nmap",
        //     "-T5", "-sn", "-PE", "-vv", "-oX", "-", target
        // ]
        
    // Log the command being executed for debugging purposes
    log.Printf("Running command: %v", cmd.Args)

    // Execute the Nmap command in Docker
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", err
    }
    return string(output), nil
}
