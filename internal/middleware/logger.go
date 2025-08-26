package middleware

import (
	"time"

	"huddle/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		userAgent := c.Request.UserAgent()

		if raw != "" {
			path = path + "?" + raw
		}

		// Log with structured logging
		logger.HTTPRequest(method, path, clientIP, statusCode, latency, userAgent)

		// Also log to console for development
		if gin.Mode() == gin.DebugMode {
			zap.L().Info("HTTP Request",
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.Int("status_code", statusCode),
				zap.Duration("latency", latency),
				zap.Int("body_size", bodySize),
				zap.String("user_agent", userAgent),
			)
		}
	}
}

// ErrorLogger logs errors with context
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("Request Error",
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.String("error", err.Error()),
				)
			}
		}
	}
}
