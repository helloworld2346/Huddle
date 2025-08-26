package auth

import (
	"huddle/pkg/logger"
	"huddle/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for auth operations
type Handler struct {
	service Service
}

// NewHandler creates a new auth handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind login request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	response, err := h.service.Login(c.Request.Context(), &req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		logger.Error("Failed to login", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, response, response.Message)
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind register request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	response, err := h.service.Register(c.Request.Context(), &req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		logger.Error("Failed to register", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, response, response.Message)
}

// Logout handles user logout
func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind logout request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	// Get access token from Authorization header
	authHeader := c.GetHeader("Authorization")
	accessToken := ""
	if authHeader != "" {
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			accessToken = tokenParts[1]
		}
	}
	
	if err := h.service.Logout(c.Request.Context(), req.RefreshToken, accessToken); err != nil {
		logger.Error("Failed to logout", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, nil, "Logout successful")
}

// RefreshToken handles token refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind refresh token request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	response, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		logger.Error("Failed to refresh token", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, response, response.Message)
}

// ForgotPassword handles forgot password request
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind forgot password request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	if err := h.service.ForgotPassword(c.Request.Context(), &req); err != nil {
		logger.Error("Failed to process forgot password", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, nil, "If the email exists, a password reset link has been sent")
}

// ResetPassword handles password reset
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind reset password request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}
	
	if err := h.service.ResetPassword(c.Request.Context(), &req); err != nil {
		logger.Error("Failed to reset password", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}
	
	utils.SuccessResponse(c, nil, "Password reset successfully")
}

// GetAuthStats handles getting auth statistics
func (h *Handler) GetAuthStats(c *gin.Context) {
	stats, err := h.service.GetAuthStats(c.Request.Context())
	if err != nil {
		logger.Error("Failed to get auth stats", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to get statistics")
		return
	}
	
	utils.SuccessResponse(c, stats, "Statistics retrieved successfully")
}