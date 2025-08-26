package conversation

import (
	"context"
)

// Repository interface defines data access methods for conversations
type Repository interface {
	// Conversations
	CreateConversation(ctx context.Context, name, convType string, createdBy uint) (*Conversation, error)
	GetConversationByID(ctx context.Context, conversationID uint) (*Conversation, error)
	GetConversationByParticipants(ctx context.Context, participantIDs []uint, convType string) (*Conversation, error)
	GetUserConversations(ctx context.Context, userID uint, limit, offset int) ([]Conversation, error)
	UpdateConversation(ctx context.Context, conversationID uint, name string) error
	DeleteConversation(ctx context.Context, conversationID uint) error

	// Participants
	AddParticipant(ctx context.Context, conversationID, userID uint, role string) (*ConversationParticipant, error)
	RemoveParticipant(ctx context.Context, conversationID, userID uint) error
	GetConversationParticipants(ctx context.Context, conversationID uint) ([]ConversationParticipant, error)
	UpdateLastReadAt(ctx context.Context, conversationID, userID uint) error
	CheckUserInConversation(ctx context.Context, conversationID, userID uint) (bool, error)
	PromoteToAdmin(ctx context.Context, conversationID, userID uint) error

	// Messages (basic operations for conversation context)
	GetLastMessage(ctx context.Context, conversationID uint) (*Message, error)
	GetUnreadCount(ctx context.Context, conversationID, userID uint) (int, error)
}

// Service interface defines business logic methods for conversations
type Service interface {
	// Conversations
	CreateConversation(ctx context.Context, userID uint, req *CreateConversationRequest) (*ConversationResponse, error)
	GetConversation(ctx context.Context, userID, conversationID uint) (*ConversationResponse, error)
	GetConversations(ctx context.Context, userID uint, limit, offset int) (*ConversationListResponse, error)
	UpdateConversation(ctx context.Context, userID, conversationID uint, req *UpdateConversationRequest) error
	DeleteConversation(ctx context.Context, userID, conversationID uint) error

	// Participants
	AddParticipant(ctx context.Context, userID, conversationID uint, req *AddParticipantRequest) error
	RemoveParticipant(ctx context.Context, userID, conversationID uint, req *RemoveParticipantRequest) error
	LeaveConversation(ctx context.Context, userID, conversationID uint, req *LeaveConversationRequest) error

	// Validation methods
	ValidateConversationAccess(ctx context.Context, userID, conversationID uint) error
	ValidateConversationAdmin(ctx context.Context, userID, conversationID uint) error
}
