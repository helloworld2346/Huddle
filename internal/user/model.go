package user

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string         `json:"username" gorm:"uniqueIndex;not null;size:20"`
	Email       string         `json:"email" gorm:"uniqueIndex;not null;size:254"`
	DisplayName string         `json:"display_name" gorm:"not null;size:50"`
	Password    string         `json:"-" gorm:"not null;size:255"` // Hidden from JSON
	Bio         string         `json:"bio" gorm:"size:500"`
	Avatar      string         `json:"avatar" gorm:"size:255"`
	IsPublic    bool           `json:"is_public" gorm:"default:true"`
	LastLogin   *time.Time     `json:"last_login"`
	LoginAttempts int          `json:"-" gorm:"default:0"`
	LockedUntil *time.Time     `json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserResponse represents user data for API responses
type UserResponse struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	Avatar      string    `json:"avatar"`
	IsPublic    bool      `json:"is_public"`
	LastLogin   *time.Time `json:"last_login"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateUserRequest represents request data for creating a user
type CreateUserRequest struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DisplayName  string `json:"display_name" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Bio          string `json:"bio"`
	IsPublic     bool   `json:"is_public"`
}

// UpdateUserRequest represents request data for updating a user
type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	Avatar      string `json:"avatar"`
	IsPublic    *bool  `json:"is_public"`
}

// ChangePasswordRequest represents request data for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// UserSearchRequest represents request data for searching users
type UserSearchRequest struct {
	Query    string `form:"q"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		Avatar:      u.Avatar,
		IsPublic:    u.IsPublic,
		LastLogin:   u.LastLogin,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// IsLocked checks if user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// IncrementLoginAttempts increments failed login attempts
func (u *User) IncrementLoginAttempts() {
	u.LoginAttempts++
	if u.LoginAttempts >= 5 {
		lockUntil := time.Now().Add(15 * time.Minute)
		u.LockedUntil = &lockUntil
	}
}

// ResetLoginAttempts resets failed login attempts
func (u *User) ResetLoginAttempts() {
	u.LoginAttempts = 0
	u.LockedUntil = nil
}

// UpdateLastLogin updates last login time
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
}
