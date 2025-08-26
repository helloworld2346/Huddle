package websocket

import (
	"huddle/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up WebSocket routes
func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	// WebSocket routes group
	ws := router.Group("/ws")
	{
		// WebSocket connection (no auth middleware - token passed via query param)
		ws.GET("/connect", handler.HandleWebSocketGin)
		
		// API endpoints (protected by auth)
		ws.GET("/users/online", middleware.AuthMiddleware(), handler.GetOnlineUsers)
		ws.GET("/users/:user_id/status", middleware.AuthMiddleware(), handler.GetUserStatus)
	}
}
