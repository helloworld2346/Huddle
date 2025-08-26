package conversation

import (
	"time"

	"huddle/internal/user"
)

// Conversation represents a chat conversation (direct or group)
type Conversation struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name"`
	Type      string    `json:"type" gorm:"not null;default:'direct';size:20"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`

	// Relations
	Creator     user.User                    `json:"creator" gorm:"foreignKey:CreatedBy"`
	Participants []ConversationParticipant   `json:"participants" gorm:"foreignKey:ConversationID"`
	Messages     []Message                   `json:"messages" gorm:"foreignKey:ConversationID"`
}

// ConversationParticipant represents a participant in a conversation
type ConversationParticipant struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID uint      `json:"conversation_id" gorm:"not null"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	Role           string    `json:"role" gorm:"not null;default:'member';size:20"`
	JoinedAt       time.Time `json:"joined_at" gorm:"default:now()"`
	LastReadAt     time.Time `json:"last_read_at" gorm:"default:now()"`

	// Relations
	Conversation Conversation `json:"conversation" gorm:"foreignKey:ConversationID"`
	User         user.User     `json:"user" gorm:"foreignKey:UserID"`
}

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
	Conversation Conversation `json:"conversation" gorm:"foreignKey:ConversationID"`
	Sender       user.User     `json:"sender" gorm:"foreignKey:SenderID"`
	ReplyTo      *Message      `json:"reply_to" gorm:"foreignKey:ReplyToID"`
	Reactions    []MessageReaction `json:"reactions" gorm:"foreignKey:MessageID"`
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

// Conversation Type Constants
const (
	ConversationTypeDirect = "direct"
	ConversationTypeGroup  = "group"
)

// Participant Role Constants
const (
	ParticipantRoleAdmin  = "admin"
	ParticipantRoleMember = "member"
)

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

// CreateConversationRequest represents request to create a conversation
type CreateConversationRequest struct {
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=direct group"`
	ParticipantIDs []uint `json:"participant_ids" binding:"required,min=1"`
}

// UpdateConversationRequest represents request to update a conversation
type UpdateConversationRequest struct {
	Name string `json:"name" binding:"required"`
}

// ConversationResponse represents conversation response
type ConversationResponse struct {
	ID           uint                    `json:"id"`
	Name         string                  `json:"name"`
	Type         string                  `json:"type"`
	CreatedBy    uint                    `json:"created_by"`
	Creator      user.UserResponse       `json:"creator"`
	Participants []ParticipantResponse   `json:"participants"`
	LastMessage  *MessageResponse        `json:"last_message,omitempty"`
	UnreadCount  int                     `json:"unread_count"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

// ParticipantResponse represents participant response
type ParticipantResponse struct {
	ID         uint            `json:"id"`
	UserID     uint            `json:"user_id"`
	User       user.UserResponse `json:"user"`
	Role       string          `json:"role"`
	JoinedAt   time.Time       `json:"joined_at"`
	LastReadAt time.Time       `json:"last_read_at"`
}

// ConversationListResponse represents conversation list response
type ConversationListResponse struct {
	Conversations []ConversationResponse `json:"conversations"`
	Total         int                    `json:"total"`
}

// AddParticipantRequest represents request to add participant
type AddParticipantRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"oneof=admin member"`
}

// RemoveParticipantRequest represents request to remove participant
type RemoveParticipantRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// LeaveConversationRequest represents request to leave conversation
type LeaveConversationRequest struct {
	NewAdminID *uint `json:"new_admin_id,omitempty"` // Optional: specify new admin when leaving
}

// MessageResponse represents message response (for conversation context)
type MessageResponse struct {
	ID          uint      `json:"id"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	SenderID    uint      `json:"sender_id"`
	Sender      user.UserResponse `json:"sender"`
	IsEdited    bool      `json:"is_edited"`
	CreatedAt   time.Time `json:"created_at"`
}
