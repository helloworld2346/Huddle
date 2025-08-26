package auth

import (
	"context"
	"time"

	"huddle/internal/user"
)

// Repository interface defines auth data access methods
type Repository interface {
	// User authentication
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	UpdateUserLoginInfo(ctx context.Context, userID uint, lastLogin *time.Time, loginAttempts int, lockedUntil *time.Time) error
	
	// Password reset
	CreatePasswordReset(ctx context.Context, reset *PasswordReset) error
	GetPasswordResetByToken(ctx context.Context, token string) (*PasswordReset, error)
	MarkPasswordResetUsed(ctx context.Context, token string) error
	DeleteExpiredPasswordResets(ctx context.Context) error
	
	// User activity
	CreateUserActivity(ctx context.Context, activity *UserActivity) error
	GetUserActivities(ctx context.Context, userID uint, limit int) ([]UserActivity, error)
	
	// Sessions
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	DeleteSession(ctx context.Context, token string) error
	DeleteUserSessions(ctx context.Context, userID uint) error
	DeleteExpiredSessions(ctx context.Context) error
	
	// Statistics
	GetAuthStats(ctx context.Context) (*AuthStats, error)
}

// Service interface defines auth business logic methods
type Service interface {
	// Authentication
	Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest, ipAddress, userAgent string) (*RegisterResponse, error)
	Logout(ctx context.Context, refreshToken string, accessToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	
	// Password management
	ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	
	// Session management
	ValidateSession(ctx context.Context, token string) (*user.User, error)
	RevokeAllSessions(ctx context.Context, userID uint) error
	
	// Activity tracking
	LogActivity(ctx context.Context, userID uint, activityType, ipAddress, userAgent string, metadata map[string]interface{}) error
	
	// Statistics
	GetAuthStats(ctx context.Context) (*AuthStats, error)
	
	// Validation
	ValidateLoginRequest(req *LoginRequest) error
	ValidateRegisterRequest(req *RegisterRequest) error
}