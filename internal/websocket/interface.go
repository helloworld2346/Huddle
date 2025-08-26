package websocket

import (
	"context"
)

// Service interface defines WebSocket business logic
type Service interface {
	// Hub management
	GetHub() *Hub
	StartHub()
	StopHub()

	// Client management
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	GetClient(clientID string) (*Client, error)
	GetOnlineUsers() []UserStatus

	// Room management
	JoinRoom(client *Client, conversationID uint) error
	LeaveRoom(client *Client, conversationID uint) error
	GetRoomClients(conversationID uint) []*Client

	// Message broadcasting
	BroadcastToRoom(conversationID uint, message *WebSocketMessage)
	BroadcastToUser(userID uint, message *WebSocketMessage)
	BroadcastToAll(message *WebSocketMessage)

	// Event handling
	HandleNewMessage(ctx context.Context, conversationID uint, messageData map[string]interface{})
	HandleMessageUpdated(ctx context.Context, conversationID uint, messageData map[string]interface{})
	HandleMessageDeleted(ctx context.Context, conversationID uint, messageID uint)
	HandleUserJoined(ctx context.Context, conversationID uint, userData map[string]interface{})
	HandleUserLeft(ctx context.Context, conversationID uint, userID uint, username string)

	// Validation
	ValidateUserInConversation(ctx context.Context, userID, conversationID uint) (bool, error)
}


