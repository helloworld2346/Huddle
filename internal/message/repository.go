package message

import (
	"context"
	"errors"
	"strings"

	"huddle/internal/database"
	"huddle/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new message repository
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// Messages

func (r *repository) CreateMessage(ctx context.Context, conversationID, senderID uint, req *CreateMessageRequest) (*Message, error) {
	message := &Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        req.Content,
		MessageType:    req.MessageType,
		FileURL:        req.FileURL,
		FileName:       req.FileName,
		FileSize:       req.FileSize,
		ReplyToID:      req.ReplyToID,
	}

	if err := r.db.WithContext(ctx).Create(message).Error; err != nil {
		logger.Error("Failed to create message", zap.Error(err))
		return nil, err
	}

	// Load relations
	if err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("ReplyTo.Sender").
		Preload("Reactions.User").
		First(message, message.ID).Error; err != nil {
		logger.Error("Failed to load message relations", zap.Error(err))
		return nil, err
	}

	logger.Info("Message created", zap.Uint("message_id", message.ID), zap.Uint("conversation_id", conversationID))
	return message, nil
}

func (r *repository) GetMessageByID(ctx context.Context, messageID uint) (*Message, error) {
	var message Message
	if err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("ReplyTo.Sender").
		Preload("Reactions.User").
		First(&message, messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("message not found")
		}
		logger.Error("Failed to get message", zap.Error(err))
		return nil, err
	}
	return &message, nil
}

func (r *repository) GetMessages(ctx context.Context, conversationID uint, limit, offset int) ([]Message, error) {
	var messages []Message
	if err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("ReplyTo.Sender").
		Preload("Reactions.User").
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		logger.Error("Failed to get messages", zap.Error(err))
		return nil, err
	}
	return messages, nil
}

func (r *repository) GetMessagesBefore(ctx context.Context, conversationID uint, beforeID uint, limit int) ([]Message, error) {
	var messages []Message
	if err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("ReplyTo.Sender").
		Preload("Reactions.User").
		Where("conversation_id = ? AND id < ?", conversationID, beforeID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error; err != nil {
		logger.Error("Failed to get messages before", zap.Error(err))
		return nil, err
	}
	return messages, nil
}

func (r *repository) UpdateMessage(ctx context.Context, messageID uint, content string) error {
	if err := r.db.WithContext(ctx).
		Model(&Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"content":    content,
			"is_edited":  true,
			"edited_at":  "NOW()",
		}).Error; err != nil {
		logger.Error("Failed to update message", zap.Error(err))
		return err
	}
	logger.Info("Message updated", zap.Uint("message_id", messageID))
	return nil
}

func (r *repository) DeleteMessage(ctx context.Context, messageID uint) error {
	if err := r.db.WithContext(ctx).Delete(&Message{}, messageID).Error; err != nil {
		logger.Error("Failed to delete message", zap.Error(err))
		return err
	}
	logger.Info("Message deleted", zap.Uint("message_id", messageID))
	return nil
}

func (r *repository) SearchMessages(ctx context.Context, conversationID uint, query string, limit, offset int) ([]Message, error) {
	var messages []Message
	searchQuery := "%" + strings.ToLower(query) + "%"
	
	if err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("ReplyTo.Sender").
		Preload("Reactions.User").
		Where("conversation_id = ? AND LOWER(content) LIKE ?", conversationID, searchQuery).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		logger.Error("Failed to search messages", zap.Error(err))
		return nil, err
	}
	return messages, nil
}

func (r *repository) GetMessageCount(ctx context.Context, conversationID uint) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&Message{}).
		Where("conversation_id = ?", conversationID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to get message count", zap.Error(err))
		return 0, err
	}
	return int(count), nil
}

// Reactions

func (r *repository) AddReaction(ctx context.Context, messageID, userID uint, reactionType string) (*MessageReaction, error) {
	reaction := &MessageReaction{
		MessageID:    messageID,
		UserID:       userID,
		ReactionType: reactionType,
	}

	if err := r.db.WithContext(ctx).Create(reaction).Error; err != nil {
		logger.Error("Failed to add reaction", zap.Error(err))
		return nil, err
	}

	// Load relations
	if err := r.db.WithContext(ctx).Preload("User").First(reaction, reaction.ID).Error; err != nil {
		logger.Error("Failed to load reaction relations", zap.Error(err))
		return nil, err
	}

	logger.Info("Reaction added", zap.Uint("message_id", messageID), zap.String("reaction_type", reactionType))
	return reaction, nil
}

func (r *repository) RemoveReaction(ctx context.Context, messageID, userID uint, reactionType string) error {
	if err := r.db.WithContext(ctx).
		Where("message_id = ? AND user_id = ? AND reaction_type = ?", messageID, userID, reactionType).
		Delete(&MessageReaction{}).Error; err != nil {
		logger.Error("Failed to remove reaction", zap.Error(err))
		return err
	}
	logger.Info("Reaction removed", zap.Uint("message_id", messageID), zap.String("reaction_type", reactionType))
	return nil
}

func (r *repository) GetMessageReactions(ctx context.Context, messageID uint) ([]MessageReaction, error) {
	var reactions []MessageReaction
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("message_id = ?", messageID).
		Find(&reactions).Error; err != nil {
		logger.Error("Failed to get message reactions", zap.Error(err))
		return nil, err
	}
	return reactions, nil
}

func (r *repository) GetUserReaction(ctx context.Context, messageID, userID uint) (*MessageReaction, error) {
	var reaction MessageReaction
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("message_id = ? AND user_id = ?", messageID, userID).
		First(&reaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No reaction found
		}
		logger.Error("Failed to get user reaction", zap.Error(err))
		return nil, err
	}
	return &reaction, nil
}

// Validation

func (r *repository) CheckUserInConversation(ctx context.Context, conversationID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Table("conversation_participants").
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check user in conversation", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func (r *repository) CheckMessageExists(ctx context.Context, messageID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&Message{}).
		Where("id = ?", messageID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check message exists", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func (r *repository) CheckMessageSender(ctx context.Context, messageID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&Message{}).
		Where("id = ? AND sender_id = ?", messageID, userID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check message sender", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}
