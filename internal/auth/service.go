package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"huddle/internal/user"
	"huddle/pkg/auth"
	"huddle/pkg/logger"
	"huddle/pkg/validation"

	"go.uber.org/zap"
)

// service implements Service interface
type service struct {
	repo Repository
}

// NewService creates a new auth service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Login handles user login
func (s *service) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	// Validate request
	if err := s.ValidateLoginRequest(req); err != nil {
		return nil, err
	}
	
	// Get user by username
	u, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// Log failed login attempt
		s.LogActivity(ctx, 0, "login_failed", ipAddress, userAgent, map[string]interface{}{
			"username": req.Username,
			"reason":   "user_not_found",
		})
		return nil, errors.New("invalid username or password")
	}
	
	// Check if account is locked
	if u.IsLocked() {
		s.LogActivity(ctx, u.ID, "login_failed", ipAddress, userAgent, map[string]interface{}{
			"reason": "account_locked",
		})
		return nil, errors.New("account is temporarily locked")
	}
	
	// Verify password
	if !auth.CheckPassword(req.Password, u.Password) {
		// Increment login attempts
		u.IncrementLoginAttempts()
		s.repo.UpdateUserLoginInfo(ctx, u.ID, nil, u.LoginAttempts, u.LockedUntil)
		
		// Log failed login attempt
		s.LogActivity(ctx, u.ID, "login_failed", ipAddress, userAgent, map[string]interface{}{
			"reason": "invalid_password",
		})
		
		return nil, errors.New("invalid username or password")
	}
	
	// Reset login attempts on successful login
	u.ResetLoginAttempts()
	u.UpdateLastLogin()
	s.repo.UpdateUserLoginInfo(ctx, u.ID, u.LastLogin, u.LoginAttempts, u.LockedUntil)
	
	// Generate tokens
	tokens, err := auth.GenerateTokenPair(u.ID, u.Username, u.Email)
	if err != nil {
		logger.Error("Failed to generate tokens", zap.Error(err))
		return nil, errors.New("failed to generate authentication tokens")
	}
	
	// Create session
	session := &Session{
		UserID:    u.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}
	
	if err := s.repo.CreateSession(ctx, session); err != nil {
		logger.Error("Failed to create session", zap.Error(err))
		return nil, errors.New("failed to create session")
	}
	
	// Log successful login
	s.LogActivity(ctx, u.ID, "login_success", ipAddress, userAgent, nil)
	
	logger.Info("User logged in successfully",
		zap.Uint("user_id", u.ID),
		zap.String("username", u.Username),
		zap.String("ip", ipAddress),
	)
	
	userResponse := u.ToResponse()
	return &LoginResponse{
		User:    &userResponse,
		Tokens:  tokens,
		Message: "Login successful",
	}, nil
}

// Register handles user registration
func (s *service) Register(ctx context.Context, req *RegisterRequest, ipAddress, userAgent string) (*RegisterResponse, error) {
	// Validate request
	if err := s.ValidateRegisterRequest(req); err != nil {
		return nil, err
	}
	
	// Check if username exists
	existingUser, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}
	
	// Check if email exists
	existingUser, err = s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}
	
	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("failed to process password")
	}
	
	// Create user
	u := &user.User{
		Username:    req.Username,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Password:    hashedPassword,
		Bio:         req.Bio,
		IsPublic:    req.IsPublic,
	}
	
	// Save user to database (using user repository)
	userRepo := user.NewRepository()
	if err := userRepo.Create(ctx, u); err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, errors.New("failed to create user")
	}
	
	// Generate tokens
	tokens, err := auth.GenerateTokenPair(u.ID, u.Username, u.Email)
	if err != nil {
		logger.Error("Failed to generate tokens", zap.Error(err))
		return nil, errors.New("failed to generate authentication tokens")
	}
	
	// Create session
	session := &Session{
		UserID:    u.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}
	
	if err := s.repo.CreateSession(ctx, session); err != nil {
		logger.Error("Failed to create session", zap.Error(err))
		return nil, errors.New("failed to create session")
	}
	
	// Log registration
	s.LogActivity(ctx, u.ID, "registration", ipAddress, userAgent, nil)
	
	logger.Info("User registered successfully",
		zap.Uint("user_id", u.ID),
		zap.String("username", u.Username),
		zap.String("email", u.Email),
		zap.String("ip", ipAddress),
	)
	
	userResponse := u.ToResponse()
	return &RegisterResponse{
		User:    &userResponse,
		Tokens:  tokens,
		Message: "Registration successful",
	}, nil
}

// Logout handles user logout
func (s *service) Logout(ctx context.Context, refreshToken string, accessToken string) error {
	// Get session to get user ID
	session, err := s.repo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		logger.Error("Failed to get session for logout", zap.Error(err))
		return errors.New("failed to logout")
	}
	
	// Blacklist refresh token in Redis
	if err := auth.BlacklistToken(ctx, refreshToken, 7*24*time.Hour); err != nil {
		logger.Error("Failed to blacklist refresh token", zap.Error(err))
	}
	
	// Blacklist access token in Redis (15 minutes)
	if accessToken != "" {
		if err := auth.BlacklistToken(ctx, accessToken, 15*time.Minute); err != nil {
			logger.Error("Failed to blacklist access token", zap.Error(err))
		}
	}
	
	// Delete session from database
	if err := s.repo.DeleteSession(ctx, refreshToken); err != nil {
		logger.Error("Failed to delete session", zap.Error(err))
		return errors.New("failed to logout")
	}
	
	// Delete user session from Redis
	if err := auth.DeleteUserSession(ctx, session.UserID); err != nil {
		logger.Error("Failed to delete user session from Redis", zap.Error(err))
	}
	
	logger.Info("User logged out successfully", zap.Uint("user_id", session.UserID))
	return nil
}

// RefreshToken handles token refresh
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	// Validate refresh token
	claims, err := auth.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	
	// Check if session exists
	session, err := s.repo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	
	// Generate new tokens
	tokens, err := auth.GenerateTokenPair(claims.UserID, claims.Username, claims.Email)
	if err != nil {
		logger.Error("Failed to generate new tokens", zap.Error(err))
		return nil, errors.New("failed to generate new tokens")
	}
	
	// Update session with new refresh token
	session.Token = tokens.RefreshToken
	session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	
	if err := s.repo.CreateSession(ctx, session); err != nil {
		logger.Error("Failed to update session", zap.Error(err))
		return nil, errors.New("failed to refresh token")
	}
	
	// Delete old session
	s.repo.DeleteSession(ctx, refreshToken)
	
	logger.Info("Token refreshed successfully", zap.Uint("user_id", claims.UserID))
	
	return &RefreshTokenResponse{
		Tokens:  tokens,
		Message: "Token refreshed successfully",
	}, nil
}

// ForgotPassword handles forgot password request
func (s *service) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	// Get user by email
	u, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if email exists or not
		logger.Info("Forgot password requested for non-existent email", zap.String("email", req.Email))
		return nil
	}
	
	// Generate reset token
	token := generateSecureToken()
	
	// Create password reset
	reset := &PasswordReset{
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
	}
	
	if err := s.repo.CreatePasswordReset(ctx, reset); err != nil {
		logger.Error("Failed to create password reset", zap.Error(err))
		return errors.New("failed to process request")
	}
	
	// Log activity
	s.LogActivity(ctx, u.ID, "forgot_password", "", "", map[string]interface{}{
		"email": req.Email,
	})
	
	logger.Info("Password reset token created", zap.Uint("user_id", u.ID), zap.String("email", req.Email))
	
	// TODO: Send email with reset token
	// For now, just log the token (in production, send via email)
	logger.Info("Password reset token", zap.String("token", token))
	
	return nil
}

// ResetPassword handles password reset
func (s *service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	// Get password reset by token
	reset, err := s.repo.GetPasswordResetByToken(ctx, req.Token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}
	
	// Validate new password
	if err := auth.ValidatePassword(req.NewPassword); err != nil {
		return err
	}
	
	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return errors.New("failed to process password")
	}
	
	// Update user password
	userRepo := user.NewRepository()
	u, err := userRepo.GetByID(ctx, reset.UserID)
	if err != nil {
		return errors.New("user not found")
	}
	
	u.Password = hashedPassword
	if err := userRepo.Update(ctx, u); err != nil {
		logger.Error("Failed to update password", zap.Error(err))
		return errors.New("failed to update password")
	}
	
	// Mark reset token as used
	if err := s.repo.MarkPasswordResetUsed(ctx, req.Token); err != nil {
		logger.Error("Failed to mark reset token as used", zap.Error(err))
	}
	
	// Revoke all sessions for security
	if err := s.RevokeAllSessions(ctx, reset.UserID); err != nil {
		logger.Error("Failed to revoke sessions", zap.Error(err))
	}
	
	// Log activity
	s.LogActivity(ctx, reset.UserID, "password_reset", "", "", nil)
	
	logger.Info("Password reset successfully", zap.Uint("user_id", reset.UserID))
	
	return nil
}

// ValidateSession validates a session token
func (s *service) ValidateSession(ctx context.Context, token string) (*user.User, error) {
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid session")
	}
	
	userRepo := user.NewRepository()
	u, err := userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	return u, nil
}

// RevokeAllSessions revokes all sessions for a user
func (s *service) RevokeAllSessions(ctx context.Context, userID uint) error {
	return s.repo.DeleteUserSessions(ctx, userID)
}

// LogActivity logs user activity
func (s *service) LogActivity(ctx context.Context, userID uint, activityType, ipAddress, userAgent string, metadata map[string]interface{}) error {
	activity := &UserActivity{
		UserID:       userID,
		ActivityType: activityType,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Metadata:     metadata,
	}
	
	return s.repo.CreateUserActivity(ctx, activity)
}

// GetAuthStats gets authentication statistics
func (s *service) GetAuthStats(ctx context.Context) (*AuthStats, error) {
	return s.repo.GetAuthStats(ctx)
}

// ValidateLoginRequest validates login request
func (s *service) ValidateLoginRequest(req *LoginRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// ValidateRegisterRequest validates register request
func (s *service) ValidateRegisterRequest(req *RegisterRequest) error {
	// Validate username
	if err := validation.ValidateUsername(req.Username); err != nil {
		return err
	}
	
	// Validate email
	if err := validation.ValidateEmail(req.Email); err != nil {
		return err
	}
	
	// Validate display name
	if err := validation.ValidateDisplayName(req.DisplayName); err != nil {
		return err
	}
	
	// Validate password
	if err := validation.ValidatePassword(req.Password); err != nil {
		return err
	}
	
	// Validate bio
	if err := validation.ValidateBio(req.Bio); err != nil {
		return err
	}
	
	return nil
}

// generateSecureToken generates a secure random token
func generateSecureToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}