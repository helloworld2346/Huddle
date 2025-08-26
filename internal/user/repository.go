package user

import (
	"context"
	"time"

	"huddle/internal/database"
	"huddle/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// repository implements Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// Create creates a new user
func (r *repository) Create(ctx context.Context, user *User) error {
	logger.Info("Creating user",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
	)
	
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID gets user by ID
func (r *repository) GetByID(ctx context.Context, id uint) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername gets user by username
func (r *repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets user by email
func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates user information
func (r *repository) Update(ctx context.Context, user *User) error {
	logger.Info("Updating user",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)
	
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user (soft delete)
func (r *repository) Delete(ctx context.Context, id uint) error {
	logger.Info("Deleting user", zap.Uint("user_id", id))
	
	return r.db.WithContext(ctx).Delete(&User{}, id).Error
}

// Search searches users by query
func (r *repository) Search(ctx context.Context, query string, page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64
	
	offset := (page - 1) * pageSize
	
	// Build search query (exclude soft deleted)
	db := r.db.WithContext(ctx).Model(&User{}).Where("deleted_at IS NULL")
	if query != "" {
		searchQuery := "%" + query + "%"
		db = db.Where("username ILIKE ? OR display_name ILIKE ? OR email ILIKE ?", 
			searchQuery, searchQuery, searchQuery)
	}
	
	// Count total
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err = db.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

// List gets paginated list of users
func (r *repository) List(ctx context.Context, page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64
	
	offset := (page - 1) * pageSize
	
	// Count total (exclude soft deleted)
	err := r.db.WithContext(ctx).Model(&User{}).Where("deleted_at IS NULL").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Get paginated results (exclude soft deleted)
	err = r.db.WithContext(ctx).Where("deleted_at IS NULL").Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

// ExistsByUsername checks if username exists
func (r *repository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail checks if email exists
func (r *repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// UpdateLoginInfo updates login-related information
func (r *repository) UpdateLoginInfo(ctx context.Context, userID uint, lastLogin *time.Time, loginAttempts int, lockedUntil *time.Time) error {
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
	
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(updates).Error
}
