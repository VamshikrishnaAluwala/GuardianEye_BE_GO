package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
)

func TestAmassHandler(c *gin.Context) {
	
	// Test simple Amass command
	cmd := exec.Command("docker", "run","caffix/amass", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "amass test failed",
			"details": err.Error(),
			"output":  string(output),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "ok",
		"amass_version": string(output),
	})
}
