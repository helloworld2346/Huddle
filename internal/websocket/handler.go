package websocket

import (
	"net/http"
	"strconv"

	"huddle/pkg/auth"
	"huddle/pkg/logger"
	"huddle/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Handler implements the Handler interface
type Handler struct {
	service Service
}

// NewHandler creates a new WebSocket handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// HandleWebSocket handles WebSocket upgrade and connection
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get user info from context (set by auth middleware)
	userID, username, err := h.getUserInfoFromContext(r)
	if err != nil {
		logger.Error("Failed to get user info from context", zap.Error(err))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection to WebSocket", zap.Error(err))
		return
	}

	// Create new client
	hub := h.service.GetHub()
	client := NewClient(hub, conn, userID, username)

	// Register client with hub
	h.service.RegisterClient(client)

	logger.Info("WebSocket client connected",
		zap.String("client_id", client.ID),
		zap.Uint("user_id", userID),
		zap.String("username", username))

	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

// GetOnlineUsers returns all online users
func (h *Handler) GetOnlineUsers(c *gin.Context) {
	users := h.service.GetOnlineUsers()
	
	utils.SuccessResponse(c, OnlineUsersResponse{
		Users: users,
	}, "Online users retrieved successfully")
}

// GetUserStatus returns status of a specific user
func (h *Handler) GetUserStatus(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	// Check if user is online
	users := h.service.GetOnlineUsers()
	var userStatus *UserStatus
	
	for _, user := range users {
		if user.UserID == uint(userID) {
			userStatus = &user
			break
		}
	}

	if userStatus == nil {
		// User is offline
		userStatus = &UserStatus{
			UserID:   uint(userID),
			IsOnline: false,
		}
	}

	utils.SuccessResponse(c, userStatus, "User status retrieved successfully")
}

// HandleWebSocketGin handles WebSocket upgrade using Gin context
func (h *Handler) HandleWebSocketGin(c *gin.Context) {
	// Get user info from context (set by auth middleware)
	userID, username, err := h.getUserInfoFromGinContext(c)
	if err != nil {
		logger.Error("Failed to get user info from Gin context", zap.Error(err))
		utils.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection to WebSocket", zap.Error(err))
		return
	}

	// Create new client
	hub := h.service.GetHub()
	client := NewClient(hub, conn, userID, username)

	// Register client with hub
	h.service.RegisterClient(client)

	logger.Info("WebSocket client connected",
		zap.String("client_id", client.ID),
		zap.Uint("user_id", userID),
		zap.String("username", username))

	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

// getUserInfoFromContext extracts user info from HTTP request context
func (h *Handler) getUserInfoFromContext(r *http.Request) (uint, string, error) {
	// This would be set by auth middleware
	// For now, we'll extract from headers or query params
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		userIDStr = r.URL.Query().Get("user_id")
	}
	
	if userIDStr == "" {
		return 0, "", ErrUnauthorized
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, "", ErrInvalidUserID
	}

	username := r.Header.Get("X-Username")
	if username == "" {
		username = r.URL.Query().Get("username")
	}

	if username == "" {
		return 0, "", ErrInvalidUsername
	}

	return uint(userID), username, nil
}

// getUserInfoFromGinContext extracts user info from Gin context or query params
func (h *Handler) getUserInfoFromGinContext(c *gin.Context) (uint, string, error) {
	// First try to get from Gin context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if exists {
		userID, ok := userIDInterface.(uint)
		if ok {
			usernameInterface, exists := c.Get("username")
			if exists {
				username, ok := usernameInterface.(string)
				if ok {
					return userID, username, nil
				}
			}
		}
	}

	// Fallback: extract from query parameters (for WebSocket connections)
	token := c.Query("token")
	if token == "" {
		return 0, "", ErrUnauthorized
	}

	// Validate token and extract user info
	claims, err := auth.ValidateToken(token)
	if err != nil {
		logger.Error("Token validation failed for WebSocket", zap.Error(err))
		return 0, "", ErrUnauthorized
	}

	// Check if token is blacklisted
	if auth.IsTokenBlacklisted(c.Request.Context(), token) {
		logger.Error("Token is blacklisted for WebSocket")
		return 0, "", ErrUnauthorized
	}

	return claims.UserID, claims.Username, nil
}

// Custom errors
var (
	ErrUnauthorized    = &WebSocketError{Code: "UNAUTHORIZED", Message: "Unauthorized"}
	ErrInvalidUserID   = &WebSocketError{Code: "INVALID_USER_ID", Message: "Invalid user ID"}
	ErrInvalidUsername = &WebSocketError{Code: "INVALID_USERNAME", Message: "Invalid username"}
)
