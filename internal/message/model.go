package message

import (
	"time"

	"huddle/internal/user"
)

// Message represents a message in a conversation
type Message struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID uint      `json:"conversation_id" gorm:"not null"`
	SenderID       uint      `json:"sender_id"`
	Content        string    `json:"content" gorm:"not null"`
	MessageType    string    `json:"message_type" gorm:"not null;default:'text';size:20"`
	FileURL        string    `json:"file_url"`
	FileName       string    `json:"file_name"`
	FileSize       int64     `json:"file_size"`
	ReplyToID      *uint     `json:"reply_to_id"`
	IsEdited       bool      `json:"is_edited" gorm:"default:false"`
	EditedAt       *time.Time `json:"edited_at"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:now()"`

	// Relations
	Conversation MessageConversation `json:"conversation" gorm:"foreignKey:ConversationID"`
	Sender       user.User            `json:"sender" gorm:"foreignKey:SenderID"`
	ReplyTo      *Message             `json:"reply_to" gorm:"foreignKey:ReplyToID"`
	Reactions    []MessageReaction    `json:"reactions" gorm:"foreignKey:MessageID"`
}

// MessageConversation represents conversation info for message context
type MessageConversation struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// MessageReaction represents a reaction to a message
type MessageReaction struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	MessageID    uint      `json:"message_id" gorm:"not null"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	ReactionType string    `json:"reaction_type" gorm:"not null;default:'like';size:20"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:now()"`

	// Relations
	Message Message   `json:"message" gorm:"foreignKey:MessageID"`
	User    user.User `json:"user" gorm:"foreignKey:UserID"`
}

// Message Type Constants
const (
	MessageTypeText   = "text"
	MessageTypeImage  = "image"
	MessageTypeFile   = "file"
	MessageTypeSystem = "system"
)

// Reaction Type Constants
const (
	ReactionTypeLike   = "like"
	ReactionTypeLove   = "love"
	ReactionTypeHaha   = "haha"
	ReactionTypeWow    = "wow"
	ReactionTypeSad    = "sad"
	ReactionTypeAngry  = "angry"
)

// DTOs for API requests/responses

// CreateMessageRequest represents request to create a message
type CreateMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	MessageType string `json:"message_type" binding:"required,oneof=text image file system"`
	FileURL     string `json:"file_url,omitempty"`
	FileName    string `json:"file_name,omitempty"`
	FileSize    int64  `json:"file_size,omitempty"`
	ReplyToID   *uint  `json:"reply_to_id,omitempty"`
}

// UpdateMessageRequest represents request to update a message
type UpdateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// MessageResponse represents message response
type MessageResponse struct {
	ID          uint                    `json:"id"`
	Content     string                  `json:"content"`
	MessageType string                  `json:"message_type"`
	SenderID    uint                    `json:"sender_id"`
	Sender      user.UserResponse       `json:"sender"`
	FileURL     string                  `json:"file_url,omitempty"`
	FileName    string                  `json:"file_name,omitempty"`
	FileSize    int64                   `json:"file_size,omitempty"`
	ReplyToID   *uint                   `json:"reply_to_id,omitempty"`
	ReplyTo     *MessageResponse        `json:"reply_to,omitempty"`
	IsEdited    bool                    `json:"is_edited"`
	EditedAt    *time.Time              `json:"edited_at,omitempty"`
	Reactions   []MessageReactionResponse `json:"reactions"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

// MessageReactionResponse represents message reaction response
type MessageReactionResponse struct {
	ID           uint      `json:"id"`
	ReactionType string    `json:"reaction_type"`
	User         user.UserResponse `json:"user"`
	CreatedAt    time.Time `json:"created_at"`
}

// MessageListResponse represents message list response
type MessageListResponse struct {
	Messages []MessageResponse `json:"messages"`
	Total    int               `json:"total"`
	HasMore  bool              `json:"has_more"`
}

// AddReactionRequest represents request to add reaction
type AddReactionRequest struct {
	ReactionType string `json:"reaction_type" binding:"required,oneof=like love haha wow sad angry"`
}

// SearchMessagesRequest represents request to search messages
type SearchMessagesRequest struct {
	Query string `json:"query" binding:"required,min=1"`
	Limit int    `json:"limit,omitempty"`
	Offset int   `json:"offset,omitempty"`
}
