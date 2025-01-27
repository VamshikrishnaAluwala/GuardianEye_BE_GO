package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "os/exec"
)

// NmapHandler handles the /nmap endpoint.
func NmapHandler(c *gin.Context) {
    domain := c.Query("domain")
    if domain == "" {
        c.JSON(http.StatusBadRequest, gin.H{"detail": "Domain parameter is required"})
        return
    }

    output, err := runNmap(domain)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "domain": domain,
        "output": output,
    })
}

// runNmap executes the nmap command and returns the output.
func runNmap(domain string) (string, error) {
    cmd := exec.Command("nmap", domain)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", err
    }
    return string(output), nil
}
