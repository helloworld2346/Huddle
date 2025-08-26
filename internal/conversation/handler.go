package conversation

import (
	"strconv"

	"huddle/pkg/logger"
	"huddle/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
}

// NewHandler creates a new conversation handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateConversation creates a new conversation
func (h *Handler) CreateConversation(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind create conversation request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	conversation, err := h.service.CreateConversation(c.Request.Context(), userID, &req)
	if err != nil {
		logger.Error("Failed to create conversation", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, conversation, "Conversation created successfully")
}

// GetConversation gets a conversation by ID
func (h *Handler) GetConversation(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	conversation, err := h.service.GetConversation(c.Request.Context(), userID, uint(conversationID))
	if err != nil {
		logger.Error("Failed to get conversation", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, conversation, "Conversation retrieved successfully")
}

// GetConversations gets user's conversations
func (h *Handler) GetConversations(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	// Parse pagination parameters
	limit := 20 // default limit
	offset := 0 // default offset
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	conversations, err := h.service.GetConversations(c.Request.Context(), userID, limit, offset)
	if err != nil {
		logger.Error("Failed to get conversations", zap.Error(err))
		utils.InternalServerErrorResponse(c, "Failed to get conversations")
		return
	}

	utils.SuccessResponse(c, conversations, "Conversations retrieved successfully")
}

// UpdateConversation updates a conversation
func (h *Handler) UpdateConversation(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	var req UpdateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind update conversation request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	err = h.service.UpdateConversation(c.Request.Context(), userID, uint(conversationID), &req)
	if err != nil {
		logger.Error("Failed to update conversation", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Conversation updated successfully")
}

// DeleteConversation deletes a conversation
func (h *Handler) DeleteConversation(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	err = h.service.DeleteConversation(c.Request.Context(), userID, uint(conversationID))
	if err != nil {
		logger.Error("Failed to delete conversation", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Conversation deleted successfully")
}

// AddParticipant adds a participant to a conversation
func (h *Handler) AddParticipant(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	var req AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind add participant request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	err = h.service.AddParticipant(c.Request.Context(), userID, uint(conversationID), &req)
	if err != nil {
		logger.Error("Failed to add participant", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Participant added successfully")
}

// RemoveParticipant removes a participant from a conversation
func (h *Handler) RemoveParticipant(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	var req RemoveParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind remove participant request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	err = h.service.RemoveParticipant(c.Request.Context(), userID, uint(conversationID), &req)
	if err != nil {
		logger.Error("Failed to remove participant", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Participant removed successfully")
}

// LeaveConversation allows a user to leave a conversation
func (h *Handler) LeaveConversation(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	var req LeaveConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If no body provided, use empty request
		req = LeaveConversationRequest{}
	}

	err = h.service.LeaveConversation(c.Request.Context(), userID, uint(conversationID), &req)
	if err != nil {
		logger.Error("Failed to leave conversation", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Left conversation successfully")
}

// Helper function to get user ID from context
func getUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}
