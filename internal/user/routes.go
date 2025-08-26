package user

import (
	"huddle/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up user routes
func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	// User routes group
	users := router.Group("/users")
	{
		// Public routes (no auth required)
		users.POST("/", handler.CreateUser)                    // Create user
		users.GET("/search", handler.SearchUsers)              // Search users
		users.GET("/", handler.ListUsers)                      // List users
		users.GET("/:id", handler.GetUserByID)                 // Get user by ID
		users.GET("/username/:username", handler.GetUserByUsername) // Get user by username
		
		// Protected routes (auth required)
		protected := users.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", handler.GetCurrentUser)           // Get current user
			protected.PUT("/me", handler.UpdateUser)               // Update current user
			protected.DELETE("/me", handler.DeleteUser)            // Delete current user
			protected.PUT("/me/password", handler.ChangePassword)  // Change password
			protected.PUT("/me/avatar", handler.UpdateAvatar)      // Update avatar
		}
	}
}
