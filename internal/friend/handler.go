package friend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"huddle/pkg/logger"
	"huddle/pkg/utils"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
}

// NewHandler creates a new friend handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Friend Request Handlers

// SendFriendRequest handles sending a friend request
func (h *Handler) SendFriendRequest(c *gin.Context) {
	var req SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind send friend request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := c.GetUint("user_id")

	response, err := h.service.SendFriendRequest(c.Request.Context(), userID, &req)
	if err != nil {
		logger.Error("Failed to send friend request", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response, "Friend request sent successfully")
}

// GetFriendRequests handles getting pending friend requests
func (h *Handler) GetFriendRequests(c *gin.Context) {
	userID := c.GetUint("user_id")

	response, err := h.service.GetFriendRequests(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get friend requests", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to get friend requests")
		return
	}

	utils.SuccessResponse(c, response, "Friend requests retrieved successfully")
}

// GetSentFriendRequests handles getting sent friend requests
func (h *Handler) GetSentFriendRequests(c *gin.Context) {
	userID := c.GetUint("user_id")

	response, err := h.service.GetSentFriendRequests(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get sent friend requests", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to get sent friend requests")
		return
	}

	utils.SuccessResponse(c, response, "Sent friend requests retrieved successfully")
}

// RespondToFriendRequest handles responding to a friend request
func (h *Handler) RespondToFriendRequest(c *gin.Context) {
	var req RespondToFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind respond to friend request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.service.RespondToFriendRequest(c.Request.Context(), userID, &req); err != nil {
		logger.Error("Failed to respond to friend request", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	action := "accepted"
	if req.Action == "reject" {
		action = "rejected"
	}

	utils.SuccessResponse(c, nil, "Friend request "+action+" successfully")
}

// CancelFriendRequest handles cancelling a friend request
func (h *Handler) CancelFriendRequest(c *gin.Context) {
	requestIDStr := c.Param("request_id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		logger.Error("Invalid request ID", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.service.CancelFriendRequest(c.Request.Context(), userID, uint(requestID)); err != nil {
		logger.Error("Failed to cancel friend request", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Friend request cancelled successfully")
}

// Friendship Handlers

// GetFriends handles getting user's friends list
func (h *Handler) GetFriends(c *gin.Context) {
	userID := c.GetUint("user_id")

	response, err := h.service.GetFriends(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get friends", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to get friends")
		return
	}

	utils.SuccessResponse(c, response, "Friends list retrieved successfully")
}

// RemoveFriend handles removing a friend
func (h *Handler) RemoveFriend(c *gin.Context) {
	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseUint(friendIDStr, 10, 32)
	if err != nil {
		logger.Error("Invalid friend ID", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid friend ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.service.RemoveFriend(c.Request.Context(), userID, uint(friendID)); err != nil {
		logger.Error("Failed to remove friend", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Friend removed successfully")
}

// CheckFriendship handles checking if two users are friends
func (h *Handler) CheckFriendship(c *gin.Context) {
	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseUint(friendIDStr, 10, 32)
	if err != nil {
		logger.Error("Invalid friend ID", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid friend ID")
		return
	}

	userID := c.GetUint("user_id")

	areFriends, err := h.service.CheckFriendship(c.Request.Context(), userID, uint(friendID))
	if err != nil {
		logger.Error("Failed to check friendship", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to check friendship")
		return
	}

	response := gin.H{
		"are_friends": areFriends,
	}

	message := "Users are not friends"
	if areFriends {
		message = "Users are friends"
	}

	utils.SuccessResponse(c, response, message)
}

// Blocked User Handlers

// BlockUser handles blocking a user
func (h *Handler) BlockUser(c *gin.Context) {
	var req BlockUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind block user request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request data")
		return
	}

	blockerID := c.GetUint("user_id")

	response, err := h.service.BlockUser(c.Request.Context(), blockerID, &req)
	if err != nil {
		logger.Error("Failed to block user", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response, "User blocked successfully")
}

// UnblockUser handles unblocking a user
func (h *Handler) UnblockUser(c *gin.Context) {
	blockedIDStr := c.Param("user_id")
	blockedID, err := strconv.ParseUint(blockedIDStr, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	blockerID := c.GetUint("user_id")

	if err := h.service.UnblockUser(c.Request.Context(), blockerID, uint(blockedID)); err != nil {
		logger.Error("Failed to unblock user", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "User unblocked successfully")
}

// GetBlockedUsers handles getting blocked users list
func (h *Handler) GetBlockedUsers(c *gin.Context) {
	blockerID := c.GetUint("user_id")

	response, err := h.service.GetBlockedUsers(c.Request.Context(), blockerID)
	if err != nil {
		logger.Error("Failed to get blocked users", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to get blocked users")
		return
	}

	utils.SuccessResponse(c, response, "Blocked users list retrieved successfully")
}

// CheckUserBlocked handles checking if a user is blocked
func (h *Handler) CheckUserBlocked(c *gin.Context) {
	blockedIDStr := c.Param("user_id")
	blockedID, err := strconv.ParseUint(blockedIDStr, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	blockerID := c.GetUint("user_id")

	isBlocked, err := h.service.CheckUserBlocked(c.Request.Context(), blockerID, uint(blockedID))
	if err != nil {
		logger.Error("Failed to check if user is blocked", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to check if user is blocked")
		return
	}

	response := gin.H{
		"is_blocked": isBlocked,
	}

	message := "User is not blocked"
	if isBlocked {
		message = "User is blocked"
	}

	utils.SuccessResponse(c, response, message)
}
