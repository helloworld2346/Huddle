package friend

import (
	"time"

	"huddle/internal/user"
)

// FriendRequest represents a friend request between users
type FriendRequest struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SenderID   uint      `json:"sender_id" gorm:"not null"`
	ReceiverID uint      `json:"receiver_id" gorm:"not null"`
	Status     string    `json:"status" gorm:"not null;default:'pending';size:20"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:now()"`
	
	// Relations
	Sender   user.User `json:"sender" gorm:"foreignKey:SenderID"`
	Receiver user.User `json:"receiver" gorm:"foreignKey:ReceiverID"`
}

// Friendship represents a friendship between two users
type Friendship struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	FriendID  uint      `json:"friend_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	
	// Relations
	User   user.User `json:"user" gorm:"foreignKey:UserID"`
	Friend user.User `json:"friend" gorm:"foreignKey:FriendID"`
}

// BlockedUser represents a blocked user relationship
type BlockedUser struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	BlockerID  uint      `json:"blocker_id" gorm:"not null"`
	BlockedID  uint      `json:"blocked_id" gorm:"not null"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:now()"`
	
	// Relations
	Blocker user.User `json:"blocker" gorm:"foreignKey:BlockerID"`
	Blocked user.User `json:"blocked" gorm:"foreignKey:BlockedID"`
}

// Request Status Constants
const (
	RequestStatusPending   = "pending"
	RequestStatusAccepted  = "accepted"
	RequestStatusRejected  = "rejected"
	RequestStatusCancelled = "cancelled"
)

// DTOs for API requests/responses

// SendFriendRequestRequest represents request to send friend request
type SendFriendRequestRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Message    string `json:"message"`
}

// RespondToFriendRequestRequest represents request to respond to friend request
type RespondToFriendRequestRequest struct {
	RequestID uint   `json:"request_id" binding:"required"`
	Action    string `json:"action" binding:"required,oneof=accept reject"`
}

// BlockUserRequest represents request to block a user
type BlockUserRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Reason string `json:"reason"`
}

// FriendRequestResponse represents friend request response
type FriendRequestResponse struct {
	ID         uint      `json:"id"`
	Sender     user.User `json:"sender"`
	Receiver   user.User `json:"receiver"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FriendshipResponse represents friendship response
type FriendshipResponse struct {
	ID        uint      `json:"id"`
	User      user.User `json:"user"`
	Friend    user.User `json:"friend"`
	CreatedAt time.Time `json:"created_at"`
}

// BlockedUserResponse represents blocked user response
type BlockedUserResponse struct {
	ID         uint      `json:"id"`
	Blocker    user.User `json:"blocker"`
	Blocked    user.User `json:"blocked"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `json:"created_at"`
}

// FriendListResponse represents friend list response
type FriendListResponse struct {
	Friends []FriendshipResponse `json:"friends"`
	Total   int                  `json:"total"`
}

// FriendRequestListResponse represents friend request list response
type FriendRequestListResponse struct {
	Requests []FriendRequestResponse `json:"requests"`
	Total    int                     `json:"total"`
}

// BlockedUserListResponse represents blocked user list response
type BlockedUserListResponse struct {
	BlockedUsers []BlockedUserResponse `json:"blocked_users"`
	Total        int                   `json:"total"`
}
