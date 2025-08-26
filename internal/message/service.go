package message

import (
	"context"
	"errors"
	"huddle/internal/user"
	"huddle/internal/websocket"
	"huddle/pkg/logger"

	"go.uber.org/zap"
)

type service struct {
	repo Repository
	wsService websocket.Service
}

// NewService creates a new message service
func NewService(repo Repository, wsService websocket.Service) Service {
	return &service{
		repo: repo,
		wsService: wsService,
	}
}

// Messages

func (s *service) CreateMessage(ctx context.Context, userID, conversationID uint, req *CreateMessageRequest) (*MessageResponse, error) {
	// Validate conversation access
	if err := s.ValidateConversationAccess(ctx, userID, conversationID); err != nil {
		return nil, err
	}

	// Validate reply message if provided
	if req.ReplyToID != nil {
		exists, err := s.repo.CheckMessageExists(ctx, *req.ReplyToID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("reply message not found")
		}
	}

	// Create message
	message, err := s.repo.CreateMessage(ctx, conversationID, userID, req)
	if err != nil {
		return nil, err
	}

	// Build response
	response := s.buildMessageResponse(ctx, message)

	// Broadcast real-time message to conversation participants
	go func() {
		// Convert MessageResponse to map for WebSocket
		messageData := map[string]interface{}{
			"id":              response.ID,
			"sender_id":       response.SenderID,
			"sender_name":     response.Sender.Username,
			"content":         response.Content,
			"message_type":    response.MessageType,
			"created_at":      response.CreatedAt,
			"updated_at":      response.UpdatedAt,
		}
		
		logger.Info("Broadcasting new message", 
			zap.Uint("conversation_id", conversationID),
			zap.Uint("message_id", response.ID),
			zap.String("content", response.Content))
		
		s.wsService.HandleNewMessage(context.Background(), conversationID, messageData)
	}()

	return response, nil
}

func (s *service) GetMessage(ctx context.Context, userID, messageID uint) (*MessageResponse, error) {
	// Get message
	message, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	// Validate conversation access
	if err := s.ValidateConversationAccess(ctx, userID, message.ConversationID); err != nil {
		return nil, err
	}

	// Build response
	return s.buildMessageResponse(ctx, message), nil
}

func (s *service) GetMessages(ctx context.Context, userID, conversationID uint, limit, offset int) (*MessageListResponse, error) {
	// Validate conversation access
	if err := s.ValidateConversationAccess(ctx, userID, conversationID); err != nil {
		return nil, err
	}

	// Set default limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	// Get messages
	messages, err := s.repo.GetMessages(ctx, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.repo.GetMessageCount(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	// Build response
	var messageResponses []MessageResponse
	for _, msg := range messages {
		messageResponses = append(messageResponses, *s.buildMessageResponse(ctx, &msg))
	}

	hasMore := offset+limit < total

	return &MessageListResponse{
		Messages: messageResponses,
		Total:    total,
		HasMore:  hasMore,
	}, nil
}

func (s *service) GetMessagesBefore(ctx context.Context, userID, conversationID uint, beforeID uint, limit int) (*MessageListResponse, error) {
	// Validate conversation access
	if err := s.ValidateConversationAccess(ctx, userID, conversationID); err != nil {
		return nil, err
	}

	// Set default limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	// Get messages
	messages, err := s.repo.GetMessagesBefore(ctx, conversationID, beforeID, limit)
	if err != nil {
		return nil, err
	}

	// Build response
	var messageResponses []MessageResponse
	for _, msg := range messages {
		messageResponses = append(messageResponses, *s.buildMessageResponse(ctx, &msg))
	}

	hasMore := len(messages) == limit

	return &MessageListResponse{
		Messages: messageResponses,
		Total:    len(messageResponses),
		HasMore:  hasMore,
	}, nil
}

func (s *service) UpdateMessage(ctx context.Context, userID, messageID uint, req *UpdateMessageRequest) error {
	// Validate message access and sender
	if err := s.ValidateMessageSender(ctx, userID, messageID); err != nil {
		return err
	}

	// Update message
	return s.repo.UpdateMessage(ctx, messageID, req.Content)
}

func (s *service) DeleteMessage(ctx context.Context, userID, messageID uint) error {
	// Validate message access and sender
	if err := s.ValidateMessageSender(ctx, userID, messageID); err != nil {
		return err
	}

	// Delete message
	return s.repo.DeleteMessage(ctx, messageID)
}

func (s *service) SearchMessages(ctx context.Context, userID, conversationID uint, req *SearchMessagesRequest) (*MessageListResponse, error) {
	// Validate conversation access
	if err := s.ValidateConversationAccess(ctx, userID, conversationID); err != nil {
		return nil, err
	}

	// Set default limit
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Search messages
	messages, err := s.repo.SearchMessages(ctx, conversationID, req.Query, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Build response
	var messageResponses []MessageResponse
	for _, msg := range messages {
		messageResponses = append(messageResponses, *s.buildMessageResponse(ctx, &msg))
	}

	return &MessageListResponse{
		Messages: messageResponses,
		Total:    len(messageResponses),
		HasMore:  len(messages) == req.Limit,
	}, nil
}

// Reactions

func (s *service) AddReaction(ctx context.Context, userID, messageID uint, req *AddReactionRequest) error {
	// Validate message access
	if err := s.ValidateMessageAccess(ctx, userID, messageID); err != nil {
		return err
	}

	// Check if reaction already exists
	existingReaction, err := s.repo.GetUserReaction(ctx, messageID, userID)
	if err != nil {
		return err
	}
	if existingReaction != nil {
		return errors.New("user already reacted to this message")
	}

	// Add reaction
	_, err = s.repo.AddReaction(ctx, messageID, userID, req.ReactionType)
	return err
}

func (s *service) RemoveReaction(ctx context.Context, userID, messageID uint, reactionType string) error {
	// Validate message access
	if err := s.ValidateMessageAccess(ctx, userID, messageID); err != nil {
		return err
	}

	// Remove reaction
	return s.repo.RemoveReaction(ctx, messageID, userID, reactionType)
}

// Validation methods

func (s *service) ValidateConversationAccess(ctx context.Context, userID, conversationID uint) error {
	isParticipant, err := s.repo.CheckUserInConversation(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.New("access denied: not a participant")
	}
	return nil
}

func (s *service) ValidateMessageAccess(ctx context.Context, userID, messageID uint) error {
	// Get message to check conversation
	message, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Validate conversation access
	return s.ValidateConversationAccess(ctx, userID, message.ConversationID)
}

func (s *service) ValidateMessageSender(ctx context.Context, userID, messageID uint) error {
	// Validate message access first
	if err := s.ValidateMessageAccess(ctx, userID, messageID); err != nil {
		return err
	}

	// Check if user is the sender
	isSender, err := s.repo.CheckMessageSender(ctx, messageID, userID)
	if err != nil {
		return err
	}
	if !isSender {
		return errors.New("access denied: not the message sender")
	}
	return nil
}

// Helper methods

func (s *service) buildMessageResponse(ctx context.Context, message *Message) *MessageResponse {
	// Build reactions response
	var reactions []MessageReactionResponse
	for _, r := range message.Reactions {
		reactions = append(reactions, MessageReactionResponse{
			ID:           r.ID,
			ReactionType: r.ReactionType,
			User: user.UserResponse{
				ID:           r.User.ID,
				Username:     r.User.Username,
				Email:        r.User.Email,
				DisplayName:  r.User.DisplayName,
				Bio:          r.User.Bio,
				Avatar:       r.User.Avatar,
				IsPublic:     r.User.IsPublic,
				LastLogin:    r.User.LastLogin,
				CreatedAt:    r.User.CreatedAt,
				UpdatedAt:    r.User.UpdatedAt,
			},
			CreatedAt: r.CreatedAt,
		})
	}

	// Build reply response if exists
	var replyTo *MessageResponse
	if message.ReplyTo != nil {
		replyTo = &MessageResponse{
			ID:          message.ReplyTo.ID,
			Content:     message.ReplyTo.Content,
			MessageType: message.ReplyTo.MessageType,
			SenderID:    message.ReplyTo.SenderID,
			Sender: user.UserResponse{
				ID:           message.ReplyTo.Sender.ID,
				Username:     message.ReplyTo.Sender.Username,
				Email:        message.ReplyTo.Sender.Email,
				DisplayName:  message.ReplyTo.Sender.DisplayName,
				Bio:          message.ReplyTo.Sender.Bio,
				Avatar:       message.ReplyTo.Sender.Avatar,
				IsPublic:     message.ReplyTo.Sender.IsPublic,
				LastLogin:    message.ReplyTo.Sender.LastLogin,
				CreatedAt:    message.ReplyTo.Sender.CreatedAt,
				UpdatedAt:    message.ReplyTo.Sender.UpdatedAt,
			},
			IsEdited:  message.ReplyTo.IsEdited,
			CreatedAt: message.ReplyTo.CreatedAt,
		}
	}

	return &MessageResponse{
		ID:          message.ID,
		Content:     message.Content,
		MessageType: message.MessageType,
		SenderID:    message.SenderID,
		Sender: user.UserResponse{
			ID:           message.Sender.ID,
			Username:     message.Sender.Username,
			Email:        message.Sender.Email,
			DisplayName:  message.Sender.DisplayName,
			Bio:          message.Sender.Bio,
			Avatar:       message.Sender.Avatar,
			IsPublic:     message.Sender.IsPublic,
			LastLogin:    message.Sender.LastLogin,
			CreatedAt:    message.Sender.CreatedAt,
			UpdatedAt:    message.Sender.UpdatedAt,
		},
		FileURL:     message.FileURL,
		FileName:    message.FileName,
		FileSize:    message.FileSize,
		ReplyToID:   message.ReplyToID,
		ReplyTo:     replyTo,
		IsEdited:    message.IsEdited,
		EditedAt:    message.EditedAt,
		Reactions:   reactions,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   message.UpdatedAt,
	}
}
