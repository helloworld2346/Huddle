package user

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, query string, page, pageSize int) ([]User, int64, error)
	List(ctx context.Context, page, pageSize int) ([]User, int64, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateLoginInfo(ctx context.Context, userID uint, lastLogin *time.Time, loginAttempts int, lockedUntil *time.Time) error
}

type Service interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*UserResponse, error)
	GetUserByUsername(ctx context.Context, username string) (*UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req *UpdateUserRequest) (*UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
	SearchUsers(ctx context.Context, req *UserSearchRequest) (*UserListResponse, error)
	ListUsers(ctx context.Context, page, pageSize int) (*UserListResponse, error)
	ChangePassword(ctx context.Context, id uint, req *ChangePasswordRequest) error
	UpdateAvatar(ctx context.Context, id uint, avatarURL string) (*UserResponse, error)
	ValidateUser(req *CreateUserRequest) error
	ValidateUpdate(req *UpdateUserRequest) error
}