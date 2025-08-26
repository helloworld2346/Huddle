package message

import (
	"huddle/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up message routes
func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	// Message routes group (all protected)
	messages := router.Group("/conversations/:id/messages")
	messages.Use(middleware.AuthMiddleware())
	{
		// Message CRUD
		messages.POST("/", handler.CreateMessage)                    // Create message
		messages.GET("/", handler.GetMessages)                       // Get messages (with pagination)
		messages.GET("/before", handler.GetMessagesBefore)           // Get messages before ID
		messages.GET("/search", handler.SearchMessages)              // Search messages
		messages.GET("/:message_id", handler.GetMessage)             // Get specific message
		messages.PUT("/:message_id", handler.UpdateMessage)          // Update message
		messages.DELETE("/:message_id", handler.DeleteMessage)       // Delete message

		// Message reactions
		messages.POST("/:message_id/reactions", handler.AddReaction)           // Add reaction
		messages.DELETE("/:message_id/reactions/:reaction_type", handler.RemoveReaction) // Remove reaction
	}
}
