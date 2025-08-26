package websocket

import (
	"encoding/json"
	"sync"
	"time"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Client events
	MessageTypeJoinConversation MessageType = "join_conversation"
	MessageTypeLeaveConversation MessageType = "leave_conversation"
	MessageTypeTyping           MessageType = "typing"
	MessageTypeStopTyping       MessageType = "stop_typing"
	MessageTypeMarkRead         MessageType = "mark_read"

	// Server events
	MessageTypeNewMessage       MessageType = "new_message"
	MessageTypeMessageUpdated   MessageType = "message_updated"
	MessageTypeMessageDeleted   MessageType = "message_deleted"
	MessageTypeUserJoined       MessageType = "user_joined"
	MessageTypeUserLeft         MessageType = "user_left"
	MessageTypeUserTyping       MessageType = "user_typing"
	MessageTypeUserStopTyping   MessageType = "user_stop_typing"
	MessageTypeUserOnline       MessageType = "user_online"
	MessageTypeUserOffline      MessageType = "user_offline"
	MessageTypeError            MessageType = "error"
	MessageTypePong             MessageType = "pong"
)

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type      MessageType          `json:"type"`
	Data      json.RawMessage      `json:"data,omitempty"`
	Timestamp time.Time            `json:"timestamp"`
	UserID    uint                 `json:"user_id,omitempty"`
	Username  string               `json:"username,omitempty"`
}

// Client represents a WebSocket client
type Client struct {
	ID         string            `json:"id"`
	UserID     uint              `json:"user_id"`
	Username   string            `json:"username"`
	Conn       interface{}       `json:"-"` // Will be *websocket.Conn
	Hub        *Hub              `json:"-"`
	Rooms      map[uint]bool     `json:"-"` // conversation_id -> bool
	Send       chan []byte       `json:"-"`
	LastPing   time.Time         `json:"last_ping"`
	IsOnline   bool              `json:"is_online"`
}

// Connection wraps the WebSocket connection
type Connection struct {
	Conn interface{} `json:"-"` // Will be *websocket.Conn
}

// Hub manages all WebSocket clients
type Hub struct {
	Clients    map[string]*Client           `json:"-"` // client_id -> client
	Rooms      map[uint]map[string]*Client  `json:"-"` // conversation_id -> clients
	Broadcast  chan *WebSocketMessage       `json:"-"`
	Register   chan *Client                 `json:"-"`
	Unregister chan *Client                 `json:"-"`
	mu         sync.RWMutex                 `json:"-"` // mutex for thread safety
	wsService  *service                     `json:"-"` // reference to service
}

// Event data structures

type JoinConversationData struct {
	ConversationID uint `json:"conversation_id"`
}

type LeaveConversationData struct {
	ConversationID uint `json:"conversation_id"`
}

type TypingData struct {
	ConversationID uint `json:"conversation_id"`
}

type MarkReadData struct {
	ConversationID uint `json:"conversation_id"`
	MessageID      uint `json:"message_id,omitempty"`
}

type NewMessageData struct {
	ConversationID uint                   `json:"conversation_id"`
	Message        map[string]interface{} `json:"message"`
}

type MessageUpdatedData struct {
	ConversationID uint                   `json:"conversation_id"`
	Message        map[string]interface{} `json:"message"`
}

type MessageDeletedData struct {
	ConversationID uint `json:"conversation_id"`
	MessageID      uint `json:"message_id"`
}

type UserJoinedData struct {
	ConversationID uint                   `json:"conversation_id"`
	User           map[string]interface{} `json:"user"`
}

type UserLeftData struct {
	ConversationID uint `json:"conversation_id"`
	UserID         uint `json:"user_id"`
	Username       string `json:"username"`
}

type UserTypingData struct {
	ConversationID uint   `json:"conversation_id"`
	UserID         uint   `json:"user_id"`
	Username       string `json:"username"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Response structures for API
type WebSocketResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type OnlineUsersResponse struct {
	Users []UserStatus `json:"users"`
}

type UserStatus struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IsOnline  bool      `json:"is_online"`
	LastSeen  time.Time `json:"last_seen,omitempty"`
}
