package message

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

// NewHandler creates a new message handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateMessage creates a new message in a conversation
func (h *Handler) CreateMessage(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind create message request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	message, err := h.service.CreateMessage(c.Request.Context(), userID, uint(conversationID), &req)
	if err != nil {
		logger.Error("Failed to create message", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, message, "Message created successfully")
}

// GetMessage gets a specific message by ID
func (h *Handler) GetMessage(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	messageID, err := strconv.ParseUint(c.Param("message_id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid message ID")
		return
	}

	message, err := h.service.GetMessage(c.Request.Context(), userID, uint(messageID))
	if err != nil {
		logger.Error("Failed to get message", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, message, "Message retrieved successfully")
}

// GetMessages gets messages from a conversation
func (h *Handler) GetMessages(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	messages, err := h.service.GetMessages(c.Request.Context(), userID, uint(conversationID), limit, offset)
	if err != nil {
		logger.Error("Failed to get messages", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, messages, "Messages retrieved successfully")
}

// GetMessagesBefore gets messages before a specific message ID
func (h *Handler) GetMessagesBefore(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	beforeID, err := strconv.ParseUint(c.Query("before"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid before ID")
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	messages, err := h.service.GetMessagesBefore(c.Request.Context(), userID, uint(conversationID), uint(beforeID), limit)
	if err != nil {
		logger.Error("Failed to get messages before", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, messages, "Messages retrieved successfully")
}

// UpdateMessage updates a message
func (h *Handler) UpdateMessage(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	messageID, err := strconv.ParseUint(c.Param("message_id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid message ID")
		return
	}

	var req UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind update message request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	err = h.service.UpdateMessage(c.Request.Context(), userID, uint(messageID), &req)
	if err != nil {
		logger.Error("Failed to update message", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Message updated successfully")
}

// DeleteMessage deletes a message
func (h *Handler) DeleteMessage(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	messageID, err := strconv.ParseUint(c.Param("message_id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid message ID")
		return
	}

	err = h.service.DeleteMessage(c.Request.Context(), userID, uint(messageID))
	if err != nil {
		logger.Error("Failed to delete message", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Message deleted successfully")
}

// SearchMessages searches messages in a conversation
func (h *Handler) SearchMessages(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid conversation ID")
		return
	}

	var req SearchMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind search messages request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	messages, err := h.service.SearchMessages(c.Request.Context(), userID, uint(conversationID), &req)
	if err != nil {
		logger.Error("Failed to search messages", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, messages, "Messages search completed")
}

// AddReaction adds a reaction to a message
func (h *Handler) AddReaction(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	messageID, err := strconv.ParseUint(c.Param("message_id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid message ID")
		return
	}

	var req AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind add reaction request", zap.Error(err))
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	err = h.service.AddReaction(c.Request.Context(), userID, uint(messageID), &req)
	if err != nil {
		logger.Error("Failed to add reaction", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Reaction added successfully")
}

// RemoveReaction removes a reaction from a message
func (h *Handler) RemoveReaction(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	messageID, err := strconv.ParseUint(c.Param("message_id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid message ID")
		return
	}

	reactionType := c.Param("reaction_type")
	if reactionType == "" {
		utils.BadRequestResponse(c, "Reaction type is required")
		return
	}

	err = h.service.RemoveReaction(c.Request.Context(), userID, uint(messageID), reactionType)
	if err != nil {
		logger.Error("Failed to remove reaction", zap.Error(err))
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Reaction removed successfully")
}

// Helper function to get user ID from context
func getUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}
