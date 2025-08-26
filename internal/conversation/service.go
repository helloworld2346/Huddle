package conversation

import (
	"context"
	"errors"
	"huddle/internal/user"
	"huddle/pkg/logger"

	"go.uber.org/zap"
)

type service struct {
	repo Repository
}

// NewService creates a new conversation service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Conversations

func (s *service) CreateConversation(ctx context.Context, userID uint, req *CreateConversationRequest) (*ConversationResponse, error) {
	// Validate request
	if err := s.validateCreateConversationRequest(req); err != nil {
		return nil, err
	}

	// Check if conversation already exists (for direct chats)
	if req.Type == ConversationTypeDirect && len(req.ParticipantIDs) == 2 {
		existingConv, err := s.repo.GetConversationByParticipants(ctx, req.ParticipantIDs, req.Type)
		if err == nil && existingConv != nil {
			// Return existing conversation
			return s.buildConversationResponse(ctx, existingConv, userID)
		}
	}

	// Create conversation
	conversation, err := s.repo.CreateConversation(ctx, req.Name, req.Type, userID)
	if err != nil {
		logger.Error("Failed to create conversation", zap.Error(err))
		return nil, err
	}

	// Add creator as admin first (ALWAYS add creator)
	logger.Info("Adding creator as admin", zap.Uint("conversation_id", conversation.ID), zap.Uint("user_id", userID))
	creatorParticipant, err := s.repo.AddParticipant(ctx, conversation.ID, userID, ParticipantRoleAdmin)
	if err != nil {
		logger.Error("Failed to add creator as participant", zap.Error(err))
		return nil, err
	}
	logger.Info("Creator added successfully", zap.Uint("conversation_id", conversation.ID), zap.Uint("user_id", userID), zap.String("role", creatorParticipant.Role))

	// Add other participants (skip creator if already in list)
	for _, participantID := range req.ParticipantIDs {
		if participantID != userID { // Skip if already added as creator
			logger.Info("Adding participant", zap.Uint("conversation_id", conversation.ID), zap.Uint("user_id", participantID))
			participant, err := s.repo.AddParticipant(ctx, conversation.ID, participantID, ParticipantRoleMember)
			if err != nil {
				logger.Error("Failed to add participant", zap.Error(err))
				return nil, err
			}
			logger.Info("Participant added successfully", zap.Uint("conversation_id", conversation.ID), zap.Uint("user_id", participantID), zap.String("role", participant.Role))
		}
	}

	// Load full conversation data
	fullConversation, err := s.repo.GetConversationByID(ctx, conversation.ID)
	if err != nil {
		logger.Error("Failed to get created conversation", zap.Error(err))
		return nil, err
	}

	return s.buildConversationResponse(ctx, fullConversation, userID)
}

func (s *service) GetConversation(ctx context.Context, userID, conversationID uint) (*ConversationResponse, error) {
	// Validate access
	if err := s.ValidateConversationAccess(ctx, userID, conversationID); err != nil {
		return nil, err
	}

	// Get conversation
	conversation, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	// Update last read timestamp
	s.repo.UpdateLastReadAt(ctx, conversationID, userID)

	return s.buildConversationResponse(ctx, conversation, userID)
}

func (s *service) GetConversations(ctx context.Context, userID uint, limit, offset int) (*ConversationListResponse, error) {
	// Get user conversations
	conversations, err := s.repo.GetUserConversations(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Build responses
	var responses []ConversationResponse
	for _, conv := range conversations {
		response, err := s.buildConversationResponse(ctx, &conv, userID)
		if err != nil {
			logger.Error("Failed to build conversation response", zap.Error(err))
			continue
		}
		responses = append(responses, *response)
	}

	return &ConversationListResponse{
		Conversations: responses,
		Total:         len(responses),
	}, nil
}

func (s *service) UpdateConversation(ctx context.Context, userID, conversationID uint, req *UpdateConversationRequest) error {
	// Validate admin access
	if err := s.ValidateConversationAdmin(ctx, userID, conversationID); err != nil {
		return err
	}

	// Update conversation
	return s.repo.UpdateConversation(ctx, conversationID, req.Name)
}

func (s *service) DeleteConversation(ctx context.Context, userID, conversationID uint) error {
	// Validate admin access
	if err := s.ValidateConversationAdmin(ctx, userID, conversationID); err != nil {
		return err
	}

	// Delete conversation
	return s.repo.DeleteConversation(ctx, conversationID)
}

// Participants

func (s *service) AddParticipant(ctx context.Context, userID, conversationID uint, req *AddParticipantRequest) error {
	// Validate admin access
	if err := s.ValidateConversationAdmin(ctx, userID, conversationID); err != nil {
		return err
	}

	// Check if user is already a participant
	isParticipant, err := s.repo.CheckUserInConversation(ctx, conversationID, req.UserID)
	if err != nil {
		return err
	}
	if isParticipant {
		return errors.New("user is already a participant")
	}

	// Add participant
	_, err = s.repo.AddParticipant(ctx, conversationID, req.UserID, req.Role)
	return err
}

func (s *service) RemoveParticipant(ctx context.Context, userID, conversationID uint, req *RemoveParticipantRequest) error {
	// Validate admin access
	if err := s.ValidateConversationAdmin(ctx, userID, conversationID); err != nil {
		return err
	}

	// Cannot remove yourself as admin
	if req.UserID == userID {
		return errors.New("cannot remove yourself as admin")
	}

	// Remove participant
	return s.repo.RemoveParticipant(ctx, conversationID, req.UserID)
}

func (s *service) LeaveConversation(ctx context.Context, userID, conversationID uint, req *LeaveConversationRequest) error {
	// Validate access
	if err := s.ValidateConversationAccess(ctx, userID, conversationID); err != nil {
		return err
	}

	// Get current participants
	participants, err := s.repo.GetConversationParticipants(ctx, conversationID)
	if err != nil {
		return err
	}

	// Find current user's role
	var isCurrentUserAdmin bool
	for _, p := range participants {
		if p.UserID == userID {
			isCurrentUserAdmin = (p.Role == ParticipantRoleAdmin)
			break
		}
	}

	// If user is admin, handle admin leave logic
	if isCurrentUserAdmin {
		// Count remaining admins (excluding current user)
		remainingAdmins := 0
		var remainingMembers []uint
		for _, p := range participants {
			if p.UserID != userID {
				if p.Role == ParticipantRoleAdmin {
					remainingAdmins++
				} else {
					remainingMembers = append(remainingMembers, p.UserID)
				}
			}
		}

		// If there are other admins, allow leave
		if remainingAdmins > 0 {
			return s.repo.RemoveParticipant(ctx, conversationID, userID)
		}

		// If no other admins, handle admin transfer
		if len(remainingMembers) == 0 {
			// No members left, delete conversation
			return s.repo.DeleteConversation(ctx, conversationID)
		} else if len(remainingMembers) == 1 {
			// Only one member left, auto-promote
			if err := s.repo.PromoteToAdmin(ctx, conversationID, remainingMembers[0]); err != nil {
				return err
			}
			logger.Info("Auto-promoted member to admin", zap.Uint("user_id", remainingMembers[0]))
		} else {
			// Multiple members, require new admin selection
			if req.NewAdminID == nil {
				return errors.New("new admin must be specified when leaving as last admin")
			}

			// Validate new admin is a member
			isValidMember := false
			for _, memberID := range remainingMembers {
				if memberID == *req.NewAdminID {
					isValidMember = true
					break
				}
			}

			if !isValidMember {
				return errors.New("specified user is not a member of this conversation")
			}

			// Promote new admin
			if err := s.repo.PromoteToAdmin(ctx, conversationID, *req.NewAdminID); err != nil {
				return err
			}
			logger.Info("Promoted member to admin", zap.Uint("user_id", *req.NewAdminID))
		}
	}

	// Leave conversation
	return s.repo.RemoveParticipant(ctx, conversationID, userID)
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

func (s *service) ValidateConversationAdmin(ctx context.Context, userID, conversationID uint) error {
	// Get participants to check role
	participants, err := s.repo.GetConversationParticipants(ctx, conversationID)
	if err != nil {
		return err
	}

	// Find user's role
	for _, participant := range participants {
		if participant.UserID == userID {
			if participant.Role == ParticipantRoleAdmin {
				return nil
			}
			return errors.New("access denied: admin role required")
		}
	}

	return errors.New("access denied: not a participant")
}

// Helper methods

func (s *service) validateCreateConversationRequest(req *CreateConversationRequest) error {
	if req.Name == "" {
		return errors.New("conversation name is required")
	}
	if req.Type != ConversationTypeDirect && req.Type != ConversationTypeGroup {
		return errors.New("invalid conversation type")
	}
	if len(req.ParticipantIDs) == 0 {
		return errors.New("at least one participant is required")
	}
	if req.Type == ConversationTypeDirect && len(req.ParticipantIDs) != 2 {
		return errors.New("direct conversations must have exactly 2 participants")
	}
	return nil
}

func (s *service) buildConversationResponse(ctx context.Context, conversation *Conversation, userID uint) (*ConversationResponse, error) {
	// Get last message
	lastMessage, err := s.repo.GetLastMessage(ctx, conversation.ID)
	if err != nil {
		logger.Error("Failed to get last message", zap.Error(err))
	}

	// Get unread count
	unreadCount, err := s.repo.GetUnreadCount(ctx, conversation.ID, userID)
	if err != nil {
		logger.Error("Failed to get unread count", zap.Error(err))
		unreadCount = 0
	}

	// Build participants response
	var participants []ParticipantResponse
	for _, p := range conversation.Participants {
		participants = append(participants, ParticipantResponse{
			ID:         p.ID,
			UserID:     p.UserID,
			User:       user.UserResponse{
				ID:           p.User.ID,
				Username:     p.User.Username,
				Email:        p.User.Email,
				DisplayName:  p.User.DisplayName,
				Bio:          p.User.Bio,
				Avatar:       p.User.Avatar,
				IsPublic:     p.User.IsPublic,
				LastLogin:    p.User.LastLogin,
				CreatedAt:    p.User.CreatedAt,
				UpdatedAt:    p.User.UpdatedAt,
			},
			Role:       p.Role,
			JoinedAt:   p.JoinedAt,
			LastReadAt: p.LastReadAt,
		})
	}

	// Build last message response
	var lastMessageResponse *MessageResponse
	if lastMessage != nil {
		lastMessageResponse = &MessageResponse{
			ID:          lastMessage.ID,
			Content:     lastMessage.Content,
			MessageType: lastMessage.MessageType,
			SenderID:    lastMessage.SenderID,
			Sender: user.UserResponse{
				ID:           lastMessage.Sender.ID,
				Username:     lastMessage.Sender.Username,
				Email:        lastMessage.Sender.Email,
				DisplayName:  lastMessage.Sender.DisplayName,
				Bio:          lastMessage.Sender.Bio,
				Avatar:       lastMessage.Sender.Avatar,
				IsPublic:     lastMessage.Sender.IsPublic,
				LastLogin:    lastMessage.Sender.LastLogin,
				CreatedAt:    lastMessage.Sender.CreatedAt,
				UpdatedAt:    lastMessage.Sender.UpdatedAt,
			},
			IsEdited:  lastMessage.IsEdited,
			CreatedAt: lastMessage.CreatedAt,
		}
	}

	return &ConversationResponse{
		ID:           conversation.ID,
		Name:         conversation.Name,
		Type:         conversation.Type,
		CreatedBy:    conversation.CreatedBy,
		Creator: user.UserResponse{
			ID:           conversation.Creator.ID,
			Username:     conversation.Creator.Username,
			Email:        conversation.Creator.Email,
			DisplayName:  conversation.Creator.DisplayName,
			Bio:          conversation.Creator.Bio,
			Avatar:       conversation.Creator.Avatar,
			IsPublic:     conversation.Creator.IsPublic,
			LastLogin:    conversation.Creator.LastLogin,
			CreatedAt:    conversation.Creator.CreatedAt,
			UpdatedAt:    conversation.Creator.UpdatedAt,
		},
		Participants: participants,
		LastMessage:  lastMessageResponse,
		UnreadCount:  unreadCount,
		CreatedAt:    conversation.CreatedAt,
		UpdatedAt:    conversation.UpdatedAt,
	}, nil
}
