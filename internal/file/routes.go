package file

import (
	"huddle/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up file-related routes
func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	// Public file access (no authentication required)
	files := router.Group("/files")
	files.GET("/:id", handler.GetFile)

	// File routes (require authentication)
	authFiles := router.Group("/files")
	authFiles.Use(middleware.AuthMiddleware())
	{
		// File upload
		authFiles.POST("/upload", handler.UploadFile)
		
		// File management
		authFiles.GET("/my", handler.ListUserFiles)
		authFiles.GET("/search", handler.SearchFiles)
		authFiles.GET("/:id/details", handler.GetFileByID)
		authFiles.PUT("/:id", handler.UpdateFile)
		authFiles.DELETE("/:id", handler.DeleteFile)
		authFiles.GET("/:id/download", handler.GetDownloadURL)
		
		// File sharing
		authFiles.POST("/share", handler.ShareFile)
		authFiles.GET("/:id/shares", handler.GetFileShares)
		authFiles.DELETE("/shares/:id", handler.DeleteFileShare)
	}

	// Conversation files (require authentication)
	conversations := router.Group("/conversations")
	conversations.Use(middleware.AuthMiddleware())
	{
		conversations.GET("/:id/files", handler.ListConversationFiles)
	}
}
