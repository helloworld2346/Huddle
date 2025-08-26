package user

import (
	"context"
	"errors"
	"math"

	"huddle/pkg/auth"
	"huddle/pkg/logger"
	"huddle/pkg/validation"

	"go.uber.org/zap"
)

// service implements Service interface
type service struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateUser creates a new user
func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	// Validate request
	if err := s.ValidateUser(req); err != nil {
		return nil, err
	}
	
	// Check if username exists
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}
	
	// Check if email exists
	exists, err = s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}
	
	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("failed to process password")
	}
	
	// Create user
	user := &User{
		Username:    req.Username,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Password:    hashedPassword,
		Bio:         req.Bio,
		IsPublic:    req.IsPublic,
	}
	
	if err := s.repo.Create(ctx, user); err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, errors.New("failed to create user")
	}
	
	logger.Info("User created successfully",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)
	
	response := user.ToResponse()
	return &response, nil
}

// GetUserByID gets user by ID
func (s *service) GetUserByID(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	response := user.ToResponse()
	return &response, nil
}

// GetUserByUsername gets user by username
func (s *service) GetUserByUsername(ctx context.Context, username string) (*UserResponse, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	
	response := user.ToResponse()
	return &response, nil
}

// UpdateUser updates user information
func (s *service) UpdateUser(ctx context.Context, id uint, req *UpdateUserRequest) (*UserResponse, error) {
	// Validate request
	if err := s.ValidateUpdate(req); err != nil {
		return nil, err
	}
	
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Update fields
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.IsPublic != nil {
		user.IsPublic = *req.IsPublic
	}
	
	// Save changes
	if err := s.repo.Update(ctx, user); err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		return nil, errors.New("failed to update user")
	}
	
	logger.Info("User updated successfully",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)
	
	response := user.ToResponse()
	return &response, nil
}

// DeleteUser deletes a user
func (s *service) DeleteUser(ctx context.Context, id uint) error {
	// Check if user exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		return errors.New("failed to delete user")
	}
	
	logger.Info("User deleted successfully", zap.Uint("user_id", id))
	return nil
}

// SearchUsers searches users by query
func (s *service) SearchUsers(ctx context.Context, req *UserSearchRequest) (*UserListResponse, error) {
	users, total, err := s.repo.Search(ctx, req.Query, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	
	// Convert to responses
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}
	
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
	
	return &UserListResponse{
		Users:      responses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ListUsers gets paginated list of users
func (s *service) ListUsers(ctx context.Context, page, pageSize int) (*UserListResponse, error) {
	users, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	
	// Convert to responses
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}
	
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	
	return &UserListResponse{
		Users:      responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ChangePassword changes user password
func (s *service) ChangePassword(ctx context.Context, id uint, req *ChangePasswordRequest) error {
	// Get user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Verify current password
	if !auth.CheckPassword(req.CurrentPassword, user.Password) {
		return errors.New("current password is incorrect")
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
	
	// Update password
	user.Password = hashedPassword
	if err := s.repo.Update(ctx, user); err != nil {
		logger.Error("Failed to update password", zap.Error(err))
		return errors.New("failed to update password")
	}
	
	logger.Info("Password changed successfully", zap.Uint("user_id", id))
	return nil
}

// UpdateAvatar updates user avatar
func (s *service) UpdateAvatar(ctx context.Context, id uint, avatarURL string) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	user.Avatar = avatarURL
	if err := s.repo.Update(ctx, user); err != nil {
		logger.Error("Failed to update avatar", zap.Error(err))
		return nil, errors.New("failed to update avatar")
	}
	
	logger.Info("Avatar updated successfully",
		zap.Uint("user_id", id),
		zap.String("avatar", avatarURL),
	)
	
	response := user.ToResponse()
	return &response, nil
}

// ValidateUser validates user data
func (s *service) ValidateUser(req *CreateUserRequest) error {
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

// ValidateUpdate validates update request
func (s *service) ValidateUpdate(req *UpdateUserRequest) error {
	// Validate display name if provided
	if req.DisplayName != "" {
		if err := validation.ValidateDisplayName(req.DisplayName); err != nil {
			return err
		}
	}
	
	// Validate bio if provided
	if req.Bio != "" {
		if err := validation.ValidateBio(req.Bio); err != nil {
			return err
		}
	}
	
	return nil
}
