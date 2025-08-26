package middleware

import (
	"strings"

	"huddle/pkg/auth"
	"huddle/pkg/logger"
	"huddle/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware validates JWT token and injects user info
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			logger.Error("Token validation failed",
				zap.String("token", tokenString[:10]+"..."),
				zap.Error(err),
			)
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Inject user info into context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("user_claims", claims)

		// Log successful authentication
		logger.Info("User authenticated",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("ip", c.ClientIP()),
		)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT token if present
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := tokenParts[1]

		// Validate token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Inject user info into context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("user_claims", claims)

		c.Next()
	}
}

// GetUserID gets user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetUsername gets username from context
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	return username.(string), true
}

// GetEmail gets email from context
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	return email.(string), true
}

// RequireAuth ensures user is authenticated
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			utils.UnauthorizedResponse(c, "Authentication required")
			c.Abort()
			return
		}

		if userID == 0 {
			utils.UnauthorizedResponse(c, "Invalid user")
			c.Abort()
			return
		}

		c.Next()
	}
}
