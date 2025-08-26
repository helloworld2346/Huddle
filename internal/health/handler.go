package health

import (
	"net/http"
	"time"

	"huddle/internal/config"
	"huddle/internal/database"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

func HealthCheck(c *gin.Context) {
	services := make(map[string]string)
	
	// Check database
	if err := database.GetDB().Raw("SELECT 1").Error; err != nil {
		services["database"] = "down"
	} else {
		services["database"] = "up"
	}
	
	// Check Redis
	redisClient := config.GetRedisClient()
	if redisClient == nil {
		services["redis"] = "down"
	} else {
		if _, err := redisClient.Ping(c.Request.Context()).Result(); err != nil {
			services["redis"] = "down"
		} else {
			services["redis"] = "up"
		}
	}
	
	// Determine overall status
	overallStatus := "healthy"
	for _, status := range services {
		if status == "down" {
			overallStatus = "unhealthy"
			break
		}
	}
	
	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Services:  services,
		Version:   "1.0.0",
	}
	
	if overallStatus == "healthy" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"timestamp": time.Now(),
	})
}
