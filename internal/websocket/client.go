package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"huddle/pkg/logger"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}



// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, userID uint, username string) *Client {
	return &Client{
		ID:       fmt.Sprintf("client_%d_%d", userID, time.Now().UnixNano()),
		UserID:   userID,
		Username: username,
		Conn:     conn,
		Hub:      hub,
		Rooms:    make(map[uint]bool),
		Send:     make(chan []byte, 256),
		LastPing: time.Now(),
		IsOnline: true,
	}
}

// getConn returns the websocket connection
func (c *Client) getConn() *websocket.Conn {
	return c.Conn.(*websocket.Conn)
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	conn := c.getConn()
	defer func() {
		c.Hub.Unregister <- c
		conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		c.LastPing = time.Now()
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket read error", zap.Error(err))
			}
			break
		}

		// Handle incoming message
		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	conn := c.getConn()
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(message []byte) {
	var wsMessage WebSocketMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		logger.Error("Failed to unmarshal WebSocket message", zap.Error(err))
		c.sendError("INVALID_MESSAGE", "Invalid message format")
		return
	}

	// Set user info
	wsMessage.UserID = c.UserID
	wsMessage.Username = c.Username
	wsMessage.Timestamp = time.Now()

	// Handle different message types
	switch wsMessage.Type {
	case MessageTypeJoinConversation:
		c.handleJoinConversation(wsMessage)
		
	case MessageTypeLeaveConversation:
		c.handleLeaveConversation(wsMessage)
		
	case MessageTypeTyping:
		c.handleTyping(wsMessage)
		
	case MessageTypeStopTyping:
		c.handleStopTyping(wsMessage)
		
	case MessageTypeMarkRead:
		c.handleMarkRead(wsMessage)
		
	default:
		logger.Warn("Unknown message type", zap.String("type", string(wsMessage.Type)))
		c.sendError("UNKNOWN_MESSAGE_TYPE", "Unknown message type")
	}
}

// handleJoinConversation handles join conversation requests
func (c *Client) handleJoinConversation(wsMessage WebSocketMessage) {
	var data JoinConversationData
	if err := json.Unmarshal(wsMessage.Data, &data); err != nil {
		c.sendError("INVALID_DATA", "Invalid join conversation data")
		return
	}

	// Validate user is in conversation
	ctx := context.Background()
	isInConversation, err := c.Hub.validateUserInConversation(ctx, c.UserID, data.ConversationID)
	if err != nil {
		logger.Error("Failed to validate user in conversation", zap.Error(err))
		c.sendError("VALIDATION_ERROR", "Failed to validate conversation access")
		return
	}

	if !isInConversation {
		c.sendError("ACCESS_DENIED", "User is not a participant in this conversation")
		return
	}

	// Join room
	c.Hub.joinRoom(c, data.ConversationID)
	
	logger.Info("Client joined conversation",
		zap.String("client_id", c.ID),
		zap.Uint("user_id", c.UserID),
		zap.Uint("conversation_id", data.ConversationID))
}

// handleLeaveConversation handles leave conversation requests
func (c *Client) handleLeaveConversation(wsMessage WebSocketMessage) {
	var data LeaveConversationData
	if err := json.Unmarshal(wsMessage.Data, &data); err != nil {
		c.sendError("INVALID_DATA", "Invalid leave conversation data")
		return
	}

	// Leave room
	c.Hub.leaveRoom(c, data.ConversationID)
	
	logger.Info("Client left conversation",
		zap.String("client_id", c.ID),
		zap.Uint("user_id", c.UserID),
		zap.Uint("conversation_id", data.ConversationID))
}

// handleTyping handles typing indicators
func (c *Client) handleTyping(wsMessage WebSocketMessage) {
	var data TypingData
	if err := json.Unmarshal(wsMessage.Data, &data); err != nil {
		c.sendError("INVALID_DATA", "Invalid typing data")
		return
	}

	// Check if user is in conversation
	if !c.Rooms[data.ConversationID] {
		c.sendError("ACCESS_DENIED", "User is not in this conversation")
		return
	}

	// Broadcast typing indicator to room
	typingMessage := &WebSocketMessage{
		Type:      MessageTypeUserTyping,
		UserID:    c.UserID,
		Username:  c.Username,
		Timestamp: time.Now(),
		Data:      wsMessage.Data,
	}

	c.Hub.Broadcast <- typingMessage
}

// handleStopTyping handles stop typing indicators
func (c *Client) handleStopTyping(wsMessage WebSocketMessage) {
	var data TypingData
	if err := json.Unmarshal(wsMessage.Data, &data); err != nil {
		c.sendError("INVALID_DATA", "Invalid stop typing data")
		return
	}

	// Check if user is in conversation
	if !c.Rooms[data.ConversationID] {
		c.sendError("ACCESS_DENIED", "User is not in this conversation")
		return
	}

	// Broadcast stop typing indicator to room
	stopTypingMessage := &WebSocketMessage{
		Type:      MessageTypeUserStopTyping,
		UserID:    c.UserID,
		Username:  c.Username,
		Timestamp: time.Now(),
		Data:      wsMessage.Data,
	}

	c.Hub.Broadcast <- stopTypingMessage
}

// handleMarkRead handles mark as read requests
func (c *Client) handleMarkRead(wsMessage WebSocketMessage) {
	var data MarkReadData
	if err := json.Unmarshal(wsMessage.Data, &data); err != nil {
		c.sendError("INVALID_DATA", "Invalid mark read data")
		return
	}

	// Check if user is in conversation
	if !c.Rooms[data.ConversationID] {
		c.sendError("ACCESS_DENIED", "User is not in this conversation")
		return
	}

	// TODO: Update last_read_at in database
	// For now, just log the action
	logger.Info("User marked conversation as read",
		zap.Uint("user_id", c.UserID),
		zap.Uint("conversation_id", data.ConversationID),
		zap.Uint("message_id", data.MessageID))
}

// sendError sends an error message to the client
func (c *Client) sendError(code, message string) {
	errorData := ErrorData{
		Code:    code,
		Message: message,
	}

	errorMessage := &WebSocketMessage{
		Type:      MessageTypeError,
		Data:      mustMarshalJSON(errorData),
		Timestamp: time.Now(),
	}

	messageBytes, err := json.Marshal(errorMessage)
	if err != nil {
		logger.Error("Failed to marshal error message", zap.Error(err))
		return
	}

	select {
	case c.Send <- messageBytes:
		// Message sent successfully
	default:
		logger.Warn("Failed to send error message, client buffer full",
			zap.String("client_id", c.ID))
	}
}


