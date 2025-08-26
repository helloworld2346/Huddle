package auth

import (
	"context"
	"time"

	"gorm.io/gorm"
	"huddle/internal/database"
	"huddle/internal/user"
	"huddle/pkg/logger"
	"go.uber.org/zap"
)

// repository implements Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new auth repository
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// GetUserByUsername gets user by username
func (r *repository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByEmail gets user by email
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateUserLoginInfo updates user login information
func (r *repository) UpdateUserLoginInfo(ctx context.Context, userID uint, lastLogin *time.Time, loginAttempts int, lockedUntil *time.Time) error {
	updates := map[string]interface{}{
		"login_attempts": loginAttempts,
	}
	
	if lastLogin != nil {
		updates["last_login"] = lastLogin
	}
	
	if lockedUntil != nil {
		updates["locked_until"] = lockedUntil
	} else {
		updates["locked_until"] = nil
	}
	
	return r.db.WithContext(ctx).Model(&user.User{}).Where("id = ?", userID).Updates(updates).Error
}

// CreatePasswordReset creates a password reset token
func (r *repository) CreatePasswordReset(ctx context.Context, reset *PasswordReset) error {
	logger.Info("Creating password reset",
		zap.Uint("user_id", reset.UserID),
		zap.String("token", reset.Token[:10]+"..."),
	)
	
	return r.db.WithContext(ctx).Create(reset).Error
}

// GetPasswordResetByToken gets password reset by token
func (r *repository) GetPasswordResetByToken(ctx context.Context, token string) (*PasswordReset, error) {
	var reset PasswordReset
	err := r.db.WithContext(ctx).Where("token = ? AND used = false AND expires_at > ?", token, time.Now()).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

// MarkPasswordResetUsed marks password reset as used
func (r *repository) MarkPasswordResetUsed(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&PasswordReset{}).Where("token = ?", token).Update("used", true).Error
}

// DeleteExpiredPasswordResets deletes expired password reset tokens
func (r *repository) DeleteExpiredPasswordResets(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ? OR used = true", time.Now()).Delete(&PasswordReset{}).Error
}

// CreateUserActivity creates user activity log
func (r *repository) CreateUserActivity(ctx context.Context, activity *UserActivity) error {
	logger.Info("Creating user activity",
		zap.Uint("user_id", activity.UserID),
		zap.String("activity_type", activity.ActivityType),
		zap.String("ip_address", activity.IPAddress),
	)
	
	return r.db.WithContext(ctx).Create(activity).Error
}

// GetUserActivities gets user activities
func (r *repository) GetUserActivities(ctx context.Context, userID uint, limit int) ([]UserActivity, error) {
	var activities []UserActivity
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&activities).Error
	if err != nil {
		return nil, err
	}
	return activities, nil
}

// CreateSession creates a user session
func (r *repository) CreateSession(ctx context.Context, session *Session) error {
	logger.Info("Creating user session",
		zap.Uint("user_id", session.UserID),
		zap.String("token", session.Token[:10]+"..."),
	)
	
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSessionByToken gets session by token
func (r *repository) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	var session Session
	err := r.db.WithContext(ctx).Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSession deletes a session
func (r *repository) DeleteSession(ctx context.Context, token string) error {
	logger.Info("Deleting session", zap.String("token", token[:10]+"..."))
	
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&Session{}).Error
}

// DeleteUserSessions deletes all sessions for a user
func (r *repository) DeleteUserSessions(ctx context.Context, userID uint) error {
	logger.Info("Deleting all sessions for user", zap.Uint("user_id", userID))
	
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&Session{}).Error
}

// DeleteExpiredSessions deletes expired sessions
func (r *repository) DeleteExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&Session{}).Error
}

// GetAuthStats gets authentication statistics
func (r *repository) GetAuthStats(ctx context.Context) (*AuthStats, error) {
	var stats AuthStats
	
	// Count total users
	err := r.db.WithContext(ctx).Model(&user.User{}).Count(&stats.TotalUsers).Error
	if err != nil {
		return nil, err
	}
	
	// Count active sessions
	err = r.db.WithContext(ctx).Model(&Session{}).Where("expires_at > ?", time.Now()).Count(&stats.ActiveSessions).Error
	if err != nil {
		return nil, err
	}
	
	// Count failed logins (last 24 hours)
	err = r.db.WithContext(ctx).Model(&UserActivity{}).Where("activity_type = ? AND created_at > ?", "login_failed", time.Now().Add(-24*time.Hour)).Count(&stats.FailedLogins).Error
	if err != nil {
		return nil, err
	}
	
	// Count password resets (last 24 hours)
	err = r.db.WithContext(ctx).Model(&PasswordReset{}).Where("created_at > ?", time.Now().Add(-24*time.Hour)).Count(&stats.PasswordResets).Error
	if err != nil {
		return nil, err
	}
	
	return &stats, nil
}