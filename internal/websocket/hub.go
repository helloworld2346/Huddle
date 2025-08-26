package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"huddle/internal/conversation"
	"huddle/pkg/logger"

	"go.uber.org/zap"
)

// NewHub creates a new WebSocket hub
func NewHub(wsService *service) *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[uint]map[string]*Client),
		Broadcast:  make(chan *WebSocketMessage, 100),
		Register:   make(chan *Client, 10),
		Unregister: make(chan *Client, 10),
		wsService:  wsService,
	}
}

// Run starts the hub goroutine
func (h *Hub) Run() {
	logger.Info("🚀 WebSocket Hub started")
	
	// Start connection health checker
	go h.connectionHealthChecker()
	
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
			
		case client := <-h.Unregister:
			h.unregisterClient(client)
			
		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// connectionHealthChecker checks for stale connections and marks users as offline
func (h *Hub) connectionHealthChecker() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()
	
	for range ticker.C {
		h.mu.Lock()
		
		now := time.Now()
		staleThreshold := 90 * time.Second // Mark as offline after 90 seconds of no ping
		
		for clientID, client := range h.Clients {
			if now.Sub(client.LastPing) > staleThreshold {
				logger.Info("Marking client as offline due to stale connection",
					zap.String("client_id", clientID),
					zap.Uint("user_id", client.UserID),
					zap.String("username", client.Username),
					zap.Duration("last_ping", now.Sub(client.LastPing)))
				
				// Mark as offline
				client.IsOnline = false
				
				// Broadcast offline status
				go h.broadcastUserStatusChange(client.UserID, client.Username, false)
				
				// Remove from hub
				delete(h.Clients, clientID)
				
				// Close connection
				close(client.Send)
			}
		}
		
		h.mu.Unlock()
	}
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.Clients[client.ID] = client
	client.IsOnline = true
	client.LastPing = time.Now()
	
	logger.Info("Client registered", 
		zap.String("client_id", client.ID),
		zap.Uint("user_id", client.UserID),
		zap.String("username", client.Username))
	
	// Broadcast online status
	go h.broadcastUserStatusChange(client.UserID, client.Username, true)
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Remove from all rooms
	for conversationID := range client.Rooms {
		if room, exists := h.Rooms[conversationID]; exists {
			delete(room, client.ID)
			
			// Remove room if empty
			if len(room) == 0 {
				delete(h.Rooms, conversationID)
			}
		}
	}
	
	// Remove from clients
	delete(h.Clients, client.ID)
	client.IsOnline = false
	
	// Close client channels
	close(client.Send)
	
	logger.Info("Client unregistered", 
		zap.String("client_id", client.ID),
		zap.Uint("user_id", client.UserID),
		zap.String("username", client.Username))
	
	// Broadcast offline status
	go h.broadcastUserStatusChange(client.UserID, client.Username, false)
}

// broadcastMessage broadcasts a message to appropriate clients
func (h *Hub) broadcastMessage(message *WebSocketMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	// Marshal message once
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logger.Error("Failed to marshal broadcast message", zap.Error(err))
		return
	}
	
	// Broadcast based on message type
	switch message.Type {
	case MessageTypeNewMessage, MessageTypeMessageUpdated, MessageTypeMessageDeleted,
		 MessageTypeUserJoined, MessageTypeUserLeft, MessageTypeUserTyping, MessageTypeUserStopTyping:
		// Extract conversation ID from data
		var conversationID uint
		if data, ok := message.Data.MarshalJSON(); ok == nil {
			var dataMap map[string]interface{}
			if json.Unmarshal(data, &dataMap) == nil {
				if convID, exists := dataMap["conversation_id"]; exists {
					if id, ok := convID.(float64); ok {
						conversationID = uint(id)
					}
				}
			}
		}
		
		if conversationID > 0 {
			h.broadcastToRoom(conversationID, messageBytes)
		}
		
	case MessageTypeUserOnline, MessageTypeUserOffline:
		// Broadcast to all clients
		h.broadcastToAll(messageBytes)
		
	default:
		// Default: broadcast to all
		h.broadcastToAll(messageBytes)
	}
}

// broadcastToRoom broadcasts message to all clients in a room
func (h *Hub) broadcastToRoom(conversationID uint, messageBytes []byte) {
	if room, exists := h.Rooms[conversationID]; exists {
		for _, client := range room {
			select {
			case client.Send <- messageBytes:
				// Message sent successfully
			default:
				// Client buffer is full, close connection
				logger.Warn("Client buffer full, closing connection",
					zap.String("client_id", client.ID),
					zap.Uint("user_id", client.UserID))
				close(client.Send)
				delete(room, client.ID)
			}
		}
	}
}

// broadcastToAll broadcasts message to all connected clients
func (h *Hub) broadcastToAll(messageBytes []byte) {
	for _, client := range h.Clients {
		select {
		case client.Send <- messageBytes:
			// Message sent successfully
		default:
			// Client buffer is full, close connection
			logger.Warn("Client buffer full, closing connection",
				zap.String("client_id", client.ID),
				zap.Uint("user_id", client.UserID))
			close(client.Send)
			delete(h.Clients, client.ID)
		}
	}
}

// joinRoom adds a client to a room
func (h *Hub) joinRoom(client *Client, conversationID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Initialize room if it doesn't exist
	if _, exists := h.Rooms[conversationID]; !exists {
		h.Rooms[conversationID] = make(map[string]*Client)
	}
	
	// Add client to room
	h.Rooms[conversationID][client.ID] = client
	client.Rooms[conversationID] = true
	
	logger.Info("Client joined room",
		zap.String("client_id", client.ID),
		zap.Uint("user_id", client.UserID),
		zap.Uint("conversation_id", conversationID))
}

// leaveRoom removes a client from a room
func (h *Hub) leaveRoom(client *Client, conversationID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if room, exists := h.Rooms[conversationID]; exists {
		delete(room, client.ID)
		delete(client.Rooms, conversationID)
		
		// Remove room if empty
		if len(room) == 0 {
			delete(h.Rooms, conversationID)
		}
		
		logger.Info("Client left room",
			zap.String("client_id", client.ID),
			zap.Uint("user_id", client.UserID),
			zap.Uint("conversation_id", conversationID))
	}
}

// getRoomClients returns all clients in a room
func (h *Hub) getRoomClients(conversationID uint) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if room, exists := h.Rooms[conversationID]; exists {
		clients := make([]*Client, 0, len(room))
		for _, client := range room {
			clients = append(clients, client)
		}
		return clients
	}
	return []*Client{}
}

// getClient returns a client by ID
func (h *Hub) getClient(clientID string) (*Client, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if client, exists := h.Clients[clientID]; exists {
		return client, nil
	}
	return nil, fmt.Errorf("client not found: %s", clientID)
}

// getOnlineUsers returns all online users
func (h *Hub) getOnlineUsers() []UserStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	users := make([]UserStatus, 0, len(h.Clients))
	for _, client := range h.Clients {
		users = append(users, UserStatus{
			UserID:   client.UserID,
			Username: client.Username,
			IsOnline: client.IsOnline,
			LastSeen: client.LastPing,
		})
	}
	return users
}

// validateUserInConversation checks if user is in conversation
func (h *Hub) validateUserInConversation(ctx context.Context, userID, conversationID uint) (bool, error) {
	// Use conversation repository to validate
	conversationRepo := conversation.NewRepository()
	return conversationRepo.CheckUserInConversation(ctx, conversationID, userID)
}

// broadcastUserStatusChange broadcasts user online/offline status to all clients
func (h *Hub) broadcastUserStatusChange(userID uint, username string, isOnline bool) {
	var messageType MessageType
	if isOnline {
		messageType = MessageTypeUserOnline
	} else {
		messageType = MessageTypeUserOffline
	}
	
	data := map[string]interface{}{
		"user_id":   userID,
		"username":  username,
		"is_online": isOnline,
		"timestamp": time.Now(),
	}
	
	message := &WebSocketMessage{
		Type:      messageType,
		Data:      mustMarshalJSON(data),
		Timestamp: time.Now(),
		UserID:    userID,
		Username:  username,
	}
	
	// Broadcast to all clients
	h.broadcastToAll(mustMarshalJSON(message))
	
	logger.Info("Broadcasted user status change",
		zap.Uint("user_id", userID),
		zap.String("username", username),
		zap.Bool("is_online", isOnline))
}
