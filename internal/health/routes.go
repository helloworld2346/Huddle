package health

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	health := r.Group("/health")
	{
		health.GET("", HealthCheck)
		health.GET("/ping", Ping)
	}
}
