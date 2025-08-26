package conversation

import (
	"github.com/gin-gonic/gin"
	"huddle/internal/middleware"
)

// SetupRoutes sets up conversation routes
func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	// Conversation routes group (all protected)
	conversations := router.Group("/conversations")
	conversations.Use(middleware.AuthMiddleware())
	{
		// Conversation CRUD
		conversations.POST("/", handler.CreateConversation)                    // Create conversation
		conversations.GET("/", handler.GetConversations)                       // Get user conversations
		conversations.GET("/:id", handler.GetConversation)                     // Get conversation by ID
		conversations.PUT("/:id", handler.UpdateConversation)                  // Update conversation
		conversations.DELETE("/:id", handler.DeleteConversation)               // Delete conversation

		// Participant management
		conversations.POST("/:id/participants", handler.AddParticipant)        // Add participant
		conversations.DELETE("/:id/participants", handler.RemoveParticipant)   // Remove participant
		conversations.POST("/:id/leave", handler.LeaveConversation)            // Leave conversation
	}
}
