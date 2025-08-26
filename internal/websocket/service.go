package websocket

import (
	"context"
	"encoding/json"
	"time"

	"huddle/internal/conversation"
	"huddle/pkg/logger"

	"go.uber.org/zap"
)

// service implements the Service interface
type service struct {
	hub *Hub
}

// NewService creates a new WebSocket service
func NewService() Service {
	s := &service{}
	s.hub = NewHub(s)
	return s
}

// GetHub returns the hub instance
func (s *service) GetHub() *Hub {
	return s.hub
}

// StartHub starts the hub goroutine
func (s *service) StartHub() {
	go s.hub.Run()
}

// StopHub stops the hub (placeholder for graceful shutdown)
func (s *service) StopHub() {
	logger.Info("🛑 WebSocket Hub stopping...")
	// TODO: Implement graceful shutdown
}

// RegisterClient registers a new client
func (s *service) RegisterClient(client *Client) {
	s.hub.Register <- client
}

// UnregisterClient unregisters a client
func (s *service) UnregisterClient(client *Client) {
	s.hub.Unregister <- client
}

// GetClient returns a client by ID
func (s *service) GetClient(clientID string) (*Client, error) {
	return s.hub.getClient(clientID)
}

// GetOnlineUsers returns all online users
func (s *service) GetOnlineUsers() []UserStatus {
	return s.hub.getOnlineUsers()
}

// JoinRoom adds a client to a room
func (s *service) JoinRoom(client *Client, conversationID uint) error {
	// Validate user is in conversation
	ctx := context.Background()
	isInConversation, err := s.ValidateUserInConversation(ctx, client.UserID, conversationID)
	if err != nil {
		logger.Error("Failed to validate user in conversation", zap.Error(err))
		return err
	}
	
	if !isInConversation {
		return ErrUserNotInConversation
	}
	
	s.hub.joinRoom(client, conversationID)
	return nil
}

// LeaveRoom removes a client from a room
func (s *service) LeaveRoom(client *Client, conversationID uint) error {
	s.hub.leaveRoom(client, conversationID)
	return nil
}

// GetRoomClients returns all clients in a room
func (s *service) GetRoomClients(conversationID uint) []*Client {
	return s.hub.getRoomClients(conversationID)
}

// BroadcastToRoom broadcasts message to all clients in a room
func (s *service) BroadcastToRoom(conversationID uint, message *WebSocketMessage) {
	s.hub.Broadcast <- message
}

// BroadcastToUser broadcasts message to a specific user
func (s *service) BroadcastToUser(userID uint, message *WebSocketMessage) {
	// Find all clients for this user
	for _, client := range s.hub.Clients {
		if client.UserID == userID {
			messageBytes, err := json.Marshal(message)
			if err != nil {
				logger.Error("Failed to marshal message", zap.Error(err))
				continue
			}
			
			select {
			case client.Send <- messageBytes:
				// Message sent successfully
			default:
				// Client buffer is full
				logger.Warn("Client buffer full", zap.String("client_id", client.ID))
			}
		}
	}
}

// BroadcastToAll broadcasts message to all connected clients
func (s *service) BroadcastToAll(message *WebSocketMessage) {
	s.hub.Broadcast <- message
}

// HandleNewMessage handles new message events
func (s *service) HandleNewMessage(ctx context.Context, conversationID uint, messageData map[string]interface{}) {
	// Add conversation_id to message data
	messageData["conversation_id"] = conversationID
	
	message := &WebSocketMessage{
		Type:      MessageTypeNewMessage,
		Data:      mustMarshalJSON(messageData),
		Timestamp: time.Now(),
	}
	
	logger.Info("Broadcasting new message", 
		zap.Uint("conversation_id", conversationID),
		zap.String("sender", messageData["sender_name"].(string)),
		zap.String("content", messageData["content"].(string)))
	
	s.BroadcastToRoom(conversationID, message)
}

// HandleMessageUpdated handles message update events
func (s *service) HandleMessageUpdated(ctx context.Context, conversationID uint, messageData map[string]interface{}) {
	message := &WebSocketMessage{
		Type:      MessageTypeMessageUpdated,
		Data:      mustMarshalJSON(messageData),
		Timestamp: time.Now(),
	}
	
	s.BroadcastToRoom(conversationID, message)
}

// HandleMessageDeleted handles message delete events
func (s *service) HandleMessageDeleted(ctx context.Context, conversationID uint, messageID uint) {
	data := map[string]interface{}{
		"conversation_id": conversationID,
		"message_id":      messageID,
	}
	
	message := &WebSocketMessage{
		Type:      MessageTypeMessageDeleted,
		Data:      mustMarshalJSON(data),
		Timestamp: time.Now(),
	}
	
	s.BroadcastToRoom(conversationID, message)
}

// HandleUserJoined handles user joined conversation events
func (s *service) HandleUserJoined(ctx context.Context, conversationID uint, userData map[string]interface{}) {
	message := &WebSocketMessage{
		Type:      MessageTypeUserJoined,
		Data:      mustMarshalJSON(userData),
		Timestamp: time.Now(),
	}
	
	s.BroadcastToRoom(conversationID, message)
}

// HandleUserLeft handles user left conversation events
func (s *service) HandleUserLeft(ctx context.Context, conversationID uint, userID uint, username string) {
	data := map[string]interface{}{
		"conversation_id": conversationID,
		"user_id":         userID,
		"username":        username,
	}
	
	message := &WebSocketMessage{
		Type:      MessageTypeUserLeft,
		Data:      mustMarshalJSON(data),
		Timestamp: time.Now(),
	}
	
	s.BroadcastToRoom(conversationID, message)
}

// ValidateUserInConversation checks if user is in conversation
func (s *service) ValidateUserInConversation(ctx context.Context, userID, conversationID uint) (bool, error) {
	conversationRepo := conversation.NewRepository()
	return conversationRepo.CheckUserInConversation(ctx, conversationID, userID)
}

// mustMarshalJSON marshals data to JSON, panics on error
func mustMarshalJSON(data interface{}) json.RawMessage {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return json.RawMessage(bytes)
}

// Custom errors
var (
	ErrUserNotInConversation = &WebSocketError{Code: "USER_NOT_IN_CONVERSATION", Message: "User is not a participant in this conversation"}
)

// WebSocketError represents a WebSocket error
type WebSocketError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *WebSocketError) Error() string {
	return e.Message
}
