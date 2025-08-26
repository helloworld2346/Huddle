package auth

import (
	"time"

	"huddle/internal/user"
	"huddle/pkg/auth"
)

// LoginRequest represents login request data
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents registration request data
type RegisterRequest struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DisplayName  string `json:"display_name" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Bio          string `json:"bio"`
	IsPublic     bool   `json:"is_public"`
}

// LoginResponse represents login response data
type LoginResponse struct {
	User         *user.UserResponse `json:"user"`
	Tokens       *auth.TokenPair    `json:"tokens"`
	Message      string             `json:"message"`
}

// RegisterResponse represents registration response data
type RegisterResponse struct {
	User         *user.UserResponse `json:"user"`
	Tokens       *auth.TokenPair    `json:"tokens"`
	Message      string             `json:"message"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	Tokens  *auth.TokenPair `json:"tokens"`
	Message string          `json:"message"`
}

// LogoutRequest represents logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// PasswordReset represents password reset token
type PasswordReset struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null;size:255"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Used      bool      `json:"used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}

// UserActivity represents user activity log
type UserActivity struct {
	ID           uint                   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint                   `json:"user_id" gorm:"not null"`
	ActivityType string                 `json:"activity_type" gorm:"not null;size:50"`
	IPAddress    string                 `json:"ip_address" gorm:"size:45"`
	UserAgent    string                 `json:"user_agent"`
	Metadata     map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt    time.Time              `json:"created_at" gorm:"default:now()"`
}

// Session represents user session
type Session struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null;size:255"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}

// AuthStats represents authentication statistics
type AuthStats struct {
	TotalUsers      int64 `json:"total_users"`
	ActiveSessions  int64 `json:"active_sessions"`
	FailedLogins    int64 `json:"failed_logins"`
	PasswordResets  int64 `json:"password_resets"`
}