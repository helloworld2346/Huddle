package auth

import (
	"github.com/gin-gonic/gin"
	"huddle/internal/middleware"
)

// SetupRoutes sets up auth routes
func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	// Auth routes group
	auth := router.Group("/auth")
	{
		// Public routes (no auth required)
		auth.POST("/login", handler.Login)                    // Login
		auth.POST("/register", handler.Register)              // Register
		auth.POST("/forgot-password", handler.ForgotPassword)  // Forgot password
		auth.POST("/reset-password", handler.ResetPassword)    // Reset password
		auth.POST("/refresh", handler.RefreshToken)           // Refresh token
		
		// Protected routes (auth required)
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/logout", handler.Logout)         // Logout
			protected.GET("/stats", handler.GetAuthStats)     // Get auth statistics
		}
	}
}