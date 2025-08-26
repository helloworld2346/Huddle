package file

import (
	"net/http"
	"strconv"

	"huddle/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// UploadFile handles file upload
// POST /api/files/upload
func (h *Handler) UploadFile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No file provided",
		})
		return
	}
	defer file.Close()

	// Parse form data
	var req UploadFileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request data",
		})
		return
	}

	// Upload file
	response, err := h.service.UploadFile(c, userID, file, header, &req)
	if err != nil {
		logger.Error("Failed to upload file", zap.Error(err), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetFile gets a file by ID (public access)
// GET /api/files/:id
func (h *Handler) GetFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	response, err := h.service.GetFile(c, uint(fileID))
	if err != nil {
		logger.Error("Failed to get file", zap.Error(err), zap.Uint("file_id", uint(fileID)))
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetFileByID gets a file by ID with access control
// GET /api/files/:id/details
func (h *Handler) GetFileByID(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	response, err := h.service.GetFileByID(c, uint(fileID), userID)
	if err != nil {
		logger.Error("Failed to get file", zap.Error(err), zap.Uint("file_id", uint(fileID)), zap.Uint("user_id", userID))
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// UpdateFile updates file metadata
// PUT /api/files/:id
func (h *Handler) UpdateFile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	var req UpdateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request data",
		})
		return
	}

	response, err := h.service.UpdateFile(c, uint(fileID), userID, &req)
	if err != nil {
		logger.Error("Failed to update file", zap.Error(err), zap.Uint("file_id", uint(fileID)), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// DeleteFile deletes a file
// DELETE /api/files/:id
func (h *Handler) DeleteFile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	if err := h.service.DeleteFile(c, uint(fileID), userID); err != nil {
		logger.Error("Failed to delete file", zap.Error(err), zap.Uint("file_id", uint(fileID)), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File deleted successfully",
	})
}

// ListUserFiles lists files owned by a user
// GET /api/files/my
func (h *Handler) ListUserFiles(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, err := h.service.ListUserFiles(c, userID, page, pageSize)
	if err != nil {
		logger.Error("Failed to list user files", zap.Error(err), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ListConversationFiles lists files in a conversation
// GET /api/conversations/:id/files
func (h *Handler) ListConversationFiles(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	conversationIDStr := c.Param("id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid conversation ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	response, err := h.service.ListConversationFiles(c, uint(conversationID), userID, page, pageSize)
	if err != nil {
		logger.Error("Failed to list conversation files", zap.Error(err), zap.Uint("conversation_id", uint(conversationID)), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// SearchFiles searches files with filters
// GET /api/files/search
func (h *Handler) SearchFiles(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	var req FileSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid query parameters",
		})
		return
	}

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	response, err := h.service.SearchFiles(c, userID, &req)
	if err != nil {
		logger.Error("Failed to search files", zap.Error(err), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// ShareFile shares a file with another user or conversation
// POST /api/files/share
func (h *Handler) ShareFile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	var req ShareFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request data",
		})
		return
	}

	response, err := h.service.ShareFile(c, userID, &req)
	if err != nil {
		logger.Error("Failed to share file", zap.Error(err), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetFileShares gets all shares for a file
// GET /api/files/:id/shares
func (h *Handler) GetFileShares(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	response, err := h.service.GetFileShares(c, uint(fileID), userID)
	if err != nil {
		logger.Error("Failed to get file shares", zap.Error(err), zap.Uint("file_id", uint(fileID)), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// DeleteFileShare deletes a file share
// DELETE /api/files/shares/:id
func (h *Handler) DeleteFileShare(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	shareIDStr := c.Param("id")
	shareID, err := strconv.ParseUint(shareIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid share ID",
		})
		return
	}

	if err := h.service.DeleteFileShare(c, uint(shareID), userID); err != nil {
		logger.Error("Failed to delete file share", zap.Error(err), zap.Uint("share_id", uint(shareID)), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File share deleted successfully",
	})
}

// GetDownloadURL generates a download URL for a file
// GET /api/files/:id/download
func (h *Handler) GetDownloadURL(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file ID",
		})
		return
	}

	url, err := h.service.GetDownloadURL(c, uint(fileID), userID)
	if err != nil {
		logger.Error("Failed to get download URL", zap.Error(err), zap.Uint("file_id", uint(fileID)), zap.Uint("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"download_url": url,
		},
	})
}

// Helper function to get user ID from context
func getUserIDFromContext(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}
