package user

import (
	"strconv"

	"huddle/internal/middleware"
	"huddle/pkg/logger"
	"huddle/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for user operations
type Handler struct {
	service Service
}

// NewHandler creates a new user handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateUser handles user creation
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind create user request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	user, err := h.service.CreateUser(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, user, "User created successfully")
}

// GetUserByID handles getting user by ID
func (h *Handler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}
	
	user, err := h.service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("Failed to get user by ID", zap.Error(err), zap.Uint64("user_id", id))
		utils.NotFoundResponse(c, "User not found")
		return
	}
	
	utils.SuccessResponse(c, user, "User retrieved successfully")
}

// GetUserByUsername handles getting user by username
func (h *Handler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.BadRequestResponse(c, "Username is required")
		return
	}
	
	user, err := h.service.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Error("Failed to get user by username", zap.Error(err), zap.String("username", username))
		utils.NotFoundResponse(c, "User not found")
		return
	}
	
	utils.SuccessResponse(c, user, "User retrieved successfully")
}

// UpdateUser handles user update
func (h *Handler) UpdateUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}
	
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind update user request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	user, err := h.service.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		logger.Error("Failed to update user", zap.Error(err), zap.Uint("user_id", userID))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, user, "User updated successfully")
}

// DeleteUser handles user deletion
func (h *Handler) DeleteUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}
	
	if err := h.service.DeleteUser(c.Request.Context(), userID); err != nil {
		logger.Error("Failed to delete user", zap.Error(err), zap.Uint("user_id", userID))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, nil, "User deleted successfully")
}

// SearchUsers handles user search
func (h *Handler) SearchUsers(c *gin.Context) {
	var req UserSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("Failed to bind search request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid query parameters")
		return
	}
	
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	
	users, err := h.service.SearchUsers(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to search users", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to search users")
		return
	}
	
	utils.SuccessResponse(c, users, "Users retrieved successfully")
}

// ListUsers handles user listing
func (h *Handler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	
	users, err := h.service.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to list users")
		return
	}
	
	utils.SuccessResponse(c, users, "Users retrieved successfully")
}

// ChangePassword handles password change
func (h *Handler) ChangePassword(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}
	
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind change password request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	if err := h.service.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		logger.Error("Failed to change password", zap.Error(err), zap.Uint("user_id", userID))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, nil, "Password changed successfully")
}

// UpdateAvatar handles avatar update
func (h *Handler) UpdateAvatar(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}
	
	// For now, we'll accept avatar URL in request body
	// Later this will be integrated with file upload
	var req struct {
		AvatarURL string `json:"avatar_url" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind avatar update request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	user, err := h.service.UpdateAvatar(c.Request.Context(), userID, req.AvatarURL)
	if err != nil {
		logger.Error("Failed to update avatar", zap.Error(err), zap.Uint("user_id", userID))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, user, "Avatar updated successfully")
}

// GetCurrentUser handles getting current user profile
func (h *Handler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}
	
	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get current user", zap.Error(err), zap.Uint("user_id", userID))
		utils.NotFoundResponse(c, "User not found")
		return
	}
	
	utils.SuccessResponse(c, user, "Current user retrieved successfully")
}