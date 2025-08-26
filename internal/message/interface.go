package message

import (
	"context"
)

// Repository interface defines data access methods for messages
type Repository interface {
	// Messages
	CreateMessage(ctx context.Context, conversationID, senderID uint, req *CreateMessageRequest) (*Message, error)
	GetMessageByID(ctx context.Context, messageID uint) (*Message, error)
	GetMessages(ctx context.Context, conversationID uint, limit, offset int) ([]Message, error)
	GetMessagesBefore(ctx context.Context, conversationID uint, beforeID uint, limit int) ([]Message, error)
	UpdateMessage(ctx context.Context, messageID uint, content string) error
	DeleteMessage(ctx context.Context, messageID uint) error
	SearchMessages(ctx context.Context, conversationID uint, query string, limit, offset int) ([]Message, error)
	GetMessageCount(ctx context.Context, conversationID uint) (int, error)

	// Reactions
	AddReaction(ctx context.Context, messageID, userID uint, reactionType string) (*MessageReaction, error)
	RemoveReaction(ctx context.Context, messageID, userID uint, reactionType string) error
	GetMessageReactions(ctx context.Context, messageID uint) ([]MessageReaction, error)
	GetUserReaction(ctx context.Context, messageID, userID uint) (*MessageReaction, error)

	// Validation
	CheckUserInConversation(ctx context.Context, conversationID, userID uint) (bool, error)
	CheckMessageExists(ctx context.Context, messageID uint) (bool, error)
	CheckMessageSender(ctx context.Context, messageID, userID uint) (bool, error)
}

// Service interface defines business logic methods for messages
type Service interface {
	// Messages
	CreateMessage(ctx context.Context, userID, conversationID uint, req *CreateMessageRequest) (*MessageResponse, error)
	GetMessage(ctx context.Context, userID, messageID uint) (*MessageResponse, error)
	GetMessages(ctx context.Context, userID, conversationID uint, limit, offset int) (*MessageListResponse, error)
	GetMessagesBefore(ctx context.Context, userID, conversationID uint, beforeID uint, limit int) (*MessageListResponse, error)
	UpdateMessage(ctx context.Context, userID, messageID uint, req *UpdateMessageRequest) error
	DeleteMessage(ctx context.Context, userID, messageID uint) error
	SearchMessages(ctx context.Context, userID, conversationID uint, req *SearchMessagesRequest) (*MessageListResponse, error)

	// Reactions
	AddReaction(ctx context.Context, userID, messageID uint, req *AddReactionRequest) error
	RemoveReaction(ctx context.Context, userID, messageID uint, reactionType string) error

	// Validation methods
	ValidateConversationAccess(ctx context.Context, userID, conversationID uint) error
	ValidateMessageAccess(ctx context.Context, userID, messageID uint) error
	ValidateMessageSender(ctx context.Context, userID, messageID uint) error
}
