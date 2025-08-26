package conversation

import (
	"context"
	"errors"
	"huddle/internal/database"
	"huddle/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new conversation repository
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// Conversations

func (r *repository) CreateConversation(ctx context.Context, name, convType string, createdBy uint) (*Conversation, error) {
	conversation := &Conversation{
		Name:      name,
		Type:      convType,
		CreatedBy: createdBy,
	}

	if err := r.db.WithContext(ctx).Create(conversation).Error; err != nil {
		logger.Error("Failed to create conversation", zap.Error(err))
		return nil, err
	}

	logger.Info("Conversation created", zap.Uint("conversation_id", conversation.ID))
	return conversation, nil
}

func (r *repository) GetConversationByID(ctx context.Context, conversationID uint) (*Conversation, error) {
	var conversation Conversation
	if err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Participants.User").
		First(&conversation, conversationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("conversation not found")
		}
		logger.Error("Failed to get conversation by ID", zap.Error(err))
		return nil, err
	}
	return &conversation, nil
}

func (r *repository) GetConversationByParticipants(ctx context.Context, participantIDs []uint, convType string) (*Conversation, error) {
	var conversation Conversation
	
	// For direct conversations, find conversation with exactly 2 participants
	if convType == ConversationTypeDirect && len(participantIDs) == 2 {
		query := r.db.WithContext(ctx).
			Joins("JOIN conversation_participants cp ON conversations.id = cp.conversation_id").
			Where("conversations.type = ?", convType).
			Where("cp.user_id IN ?", participantIDs).
			Group("conversations.id").
			Having("COUNT(DISTINCT cp.user_id) = ?", len(participantIDs)).
			Preload("Creator").
			Preload("Participants.User")
		
		if err := query.First(&conversation).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("conversation not found")
			}
			logger.Error("Failed to get conversation by participants", zap.Error(err))
			return nil, err
		}
		return &conversation, nil
	}
	
	return nil, errors.New("conversation not found")
}

func (r *repository) GetUserConversations(ctx context.Context, userID uint, limit, offset int) ([]Conversation, error) {
	var conversations []Conversation
	
	query := r.db.WithContext(ctx).
		Joins("JOIN conversation_participants cp ON conversations.id = cp.conversation_id").
		Where("cp.user_id = ?", userID).
		Preload("Creator").
		Preload("Participants.User").
		Order("conversations.updated_at DESC").
		Limit(limit).
		Offset(offset)
	
	if err := query.Find(&conversations).Error; err != nil {
		logger.Error("Failed to get user conversations", zap.Error(err))
		return nil, err
	}
	
	return conversations, nil
}

func (r *repository) UpdateConversation(ctx context.Context, conversationID uint, name string) error {
	if err := r.db.WithContext(ctx).Model(&Conversation{}).Where("id = ?", conversationID).Update("name", name).Error; err != nil {
		logger.Error("Failed to update conversation", zap.Error(err))
		return err
	}
	logger.Info("Conversation updated", zap.Uint("conversation_id", conversationID))
	return nil
}

func (r *repository) DeleteConversation(ctx context.Context, conversationID uint) error {
	if err := r.db.WithContext(ctx).Delete(&Conversation{}, conversationID).Error; err != nil {
		logger.Error("Failed to delete conversation", zap.Error(err))
		return err
	}
	logger.Info("Conversation deleted", zap.Uint("conversation_id", conversationID))
	return nil
}

// Participants

func (r *repository) AddParticipant(ctx context.Context, conversationID, userID uint, role string) (*ConversationParticipant, error) {
	participant := &ConversationParticipant{
		ConversationID: conversationID,
		UserID:         userID,
		Role:           role,
	}

	if err := r.db.WithContext(ctx).Create(participant).Error; err != nil {
		logger.Error("Failed to add participant", zap.Error(err))
		return nil, err
	}

	// Load relations
	if err := r.db.WithContext(ctx).Preload("User").First(participant, participant.ID).Error; err != nil {
		logger.Error("Failed to load participant relations", zap.Error(err))
		return nil, err
	}

	logger.Info("Participant added", zap.Uint("conversation_id", conversationID), zap.Uint("user_id", userID))
	return participant, nil
}

func (r *repository) RemoveParticipant(ctx context.Context, conversationID, userID uint) error {
	if err := r.db.WithContext(ctx).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Delete(&ConversationParticipant{}).Error; err != nil {
		logger.Error("Failed to remove participant", zap.Error(err))
		return err
	}
	logger.Info("Participant removed", zap.Uint("conversation_id", conversationID), zap.Uint("user_id", userID))
	return nil
}

func (r *repository) GetConversationParticipants(ctx context.Context, conversationID uint) ([]ConversationParticipant, error) {
	var participants []ConversationParticipant
	if err := r.db.WithContext(ctx).Preload("User").Where("conversation_id = ?", conversationID).Find(&participants).Error; err != nil {
		logger.Error("Failed to get conversation participants", zap.Error(err))
		return nil, err
	}
	return participants, nil
}

func (r *repository) UpdateLastReadAt(ctx context.Context, conversationID, userID uint) error {
	if err := r.db.WithContext(ctx).Model(&ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("last_read_at", "NOW()").Error; err != nil {
		logger.Error("Failed to update last read at", zap.Error(err))
		return err
	}
	return nil
}

func (r *repository) CheckUserInConversation(ctx context.Context, conversationID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check user in conversation", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func (r *repository) PromoteToAdmin(ctx context.Context, conversationID, userID uint) error {
	if err := r.db.WithContext(ctx).Model(&ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("role", ParticipantRoleAdmin).Error; err != nil {
		logger.Error("Failed to promote user to admin", zap.Error(err))
		return err
	}
	logger.Info("User promoted to admin", zap.Uint("conversation_id", conversationID), zap.Uint("user_id", userID))
	return nil
}

// Messages (basic operations for conversation context)

func (r *repository) GetLastMessage(ctx context.Context, conversationID uint) (*Message, error) {
	var message Message
	if err := r.db.WithContext(ctx).
		Preload("Sender").
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		First(&message).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No messages yet
		}
		logger.Error("Failed to get last message", zap.Error(err))
		return nil, err
	}
	return &message, nil
}

func (r *repository) GetUnreadCount(ctx context.Context, conversationID, userID uint) (int, error) {
	var count int64
	
	// Get user's last read timestamp
	var participant ConversationParticipant
	if err := r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		First(&participant).Error; err != nil {
		logger.Error("Failed to get participant for unread count", zap.Error(err))
		return 0, err
	}
	
	// Count messages after last read
	if err := r.db.WithContext(ctx).Model(&Message{}).
		Where("conversation_id = ? AND created_at > ? AND sender_id != ?", 
			conversationID, participant.LastReadAt, userID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to get unread count", zap.Error(err))
		return 0, err
	}
	
	return int(count), nil
}
