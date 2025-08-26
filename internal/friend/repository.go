package friend

import (
	"context"
	"errors"

	"huddle/internal/database"
	"huddle/internal/user"
	"huddle/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new friend repository
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// Friend Requests

func (r *repository) CreateFriendRequest(ctx context.Context, senderID, receiverID uint, message string) (*FriendRequest, error) {
	friendRequest := &FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     RequestStatusPending,
		Message:    message,
	}

	if err := r.db.WithContext(ctx).Create(friendRequest).Error; err != nil {
		logger.Error("Failed to create friend request", zap.Error(err))
		return nil, err
	}

	// Load relations
	if err := r.db.WithContext(ctx).Preload("Sender").Preload("Receiver").First(friendRequest, friendRequest.ID).Error; err != nil {
		logger.Error("Failed to load friend request relations", zap.Error(err))
		return nil, err
	}

	logger.Info("Friend request created", zap.Uint("request_id", friendRequest.ID))
	return friendRequest, nil
}

func (r *repository) GetFriendRequestByID(ctx context.Context, requestID uint) (*FriendRequest, error) {
	var friendRequest FriendRequest
	if err := r.db.WithContext(ctx).Preload("Sender").Preload("Receiver").First(&friendRequest, requestID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("friend request not found")
		}
		logger.Error("Failed to get friend request by ID", zap.Error(err))
		return nil, err
	}
	return &friendRequest, nil
}

func (r *repository) GetFriendRequestByUsers(ctx context.Context, senderID, receiverID uint) (*FriendRequest, error) {
	var friendRequest FriendRequest
	if err := r.db.WithContext(ctx).Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", 
			senderID, receiverID, receiverID, senderID).
		First(&friendRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("friend request not found")
		}
		logger.Error("Failed to get friend request by users", zap.Error(err))
		return nil, err
	}
	return &friendRequest, nil
}

func (r *repository) GetPendingFriendRequests(ctx context.Context, receiverID uint) ([]FriendRequest, error) {
	var requests []FriendRequest
	if err := r.db.WithContext(ctx).Preload("Sender").Preload("Receiver").
		Where("receiver_id = ? AND status = ?", receiverID, RequestStatusPending).
		Order("created_at DESC").
		Find(&requests).Error; err != nil {
		logger.Error("Failed to get pending friend requests", zap.Error(err))
		return nil, err
	}
	return requests, nil
}

func (r *repository) GetSentFriendRequests(ctx context.Context, senderID uint) ([]FriendRequest, error) {
	var requests []FriendRequest
	if err := r.db.WithContext(ctx).Preload("Sender").Preload("Receiver").
		Where("sender_id = ?", senderID).
		Order("created_at DESC").
		Find(&requests).Error; err != nil {
		logger.Error("Failed to get sent friend requests", zap.Error(err))
		return nil, err
	}
	return requests, nil
}

func (r *repository) UpdateFriendRequestStatus(ctx context.Context, requestID uint, status string) error {
	if err := r.db.WithContext(ctx).Model(&FriendRequest{}).Where("id = ?", requestID).Update("status", status).Error; err != nil {
		logger.Error("Failed to update friend request status", zap.Error(err))
		return err
	}
	logger.Info("Friend request status updated", zap.Uint("request_id", requestID), zap.String("status", status))
	return nil
}

func (r *repository) DeleteFriendRequest(ctx context.Context, requestID uint) error {
	if err := r.db.WithContext(ctx).Delete(&FriendRequest{}, requestID).Error; err != nil {
		logger.Error("Failed to delete friend request", zap.Error(err))
		return err
	}
	logger.Info("Friend request deleted", zap.Uint("request_id", requestID))
	return nil
}

// Friendships

func (r *repository) CreateFriendship(ctx context.Context, userID, friendID uint) (*Friendship, error) {
	// Create bidirectional friendship
	friendship1 := &Friendship{
		UserID:   userID,
		FriendID: friendID,
	}
	friendship2 := &Friendship{
		UserID:   friendID,
		FriendID: userID,
	}

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(friendship1).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to create friendship 1", zap.Error(err))
		return nil, err
	}

	if err := tx.Create(friendship2).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to create friendship 2", zap.Error(err))
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit friendship transaction", zap.Error(err))
		return nil, err
	}

	// Load relations
	if err := r.db.WithContext(ctx).Preload("User").Preload("Friend").First(friendship1, friendship1.ID).Error; err != nil {
		logger.Error("Failed to load friendship relations", zap.Error(err))
		return nil, err
	}

	logger.Info("Friendship created", zap.Uint("user_id", userID), zap.Uint("friend_id", friendID))
	return friendship1, nil
}

func (r *repository) GetFriendship(ctx context.Context, userID, friendID uint) (*Friendship, error) {
	var friendship Friendship
	if err := r.db.WithContext(ctx).Preload("User").Preload("Friend").
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		First(&friendship).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("friendship not found")
		}
		logger.Error("Failed to get friendship", zap.Error(err))
		return nil, err
	}
	return &friendship, nil
}

func (r *repository) GetUserFriends(ctx context.Context, userID uint) ([]Friendship, error) {
	var friendships []Friendship
	if err := r.db.WithContext(ctx).Preload("User").Preload("Friend").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&friendships).Error; err != nil {
		logger.Error("Failed to get user friends", zap.Error(err))
		return nil, err
	}
	return friendships, nil
}

func (r *repository) DeleteFriendship(ctx context.Context, userID, friendID uint) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete both directions
	if err := tx.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", 
		userID, friendID, friendID, userID).Delete(&Friendship{}).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to delete friendship", zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit friendship deletion", zap.Error(err))
		return err
	}

	logger.Info("Friendship deleted", zap.Uint("user_id", userID), zap.Uint("friend_id", friendID))
	return nil
}

func (r *repository) CheckFriendshipExists(ctx context.Context, userID, friendID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Friendship{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check friendship exists", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

// Blocked Users

func (r *repository) BlockUser(ctx context.Context, blockerID, blockedID uint, reason string) (*BlockedUser, error) {
	blockedUser := &BlockedUser{
		BlockerID: blockerID,
		BlockedID: blockedID,
		Reason:    reason,
	}

	if err := r.db.WithContext(ctx).Create(blockedUser).Error; err != nil {
		logger.Error("Failed to block user", zap.Error(err))
		return nil, err
	}

	// Load relations
	if err := r.db.WithContext(ctx).Preload("Blocker").Preload("Blocked").First(blockedUser, blockedUser.ID).Error; err != nil {
		logger.Error("Failed to load blocked user relations", zap.Error(err))
		return nil, err
	}

	logger.Info("User blocked", zap.Uint("blocker_id", blockerID), zap.Uint("blocked_id", blockedID))
	return blockedUser, nil
}

func (r *repository) UnblockUser(ctx context.Context, blockerID, blockedID uint) error {
	if err := r.db.WithContext(ctx).Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).Delete(&BlockedUser{}).Error; err != nil {
		logger.Error("Failed to unblock user", zap.Error(err))
		return err
	}
	logger.Info("User unblocked", zap.Uint("blocker_id", blockerID), zap.Uint("blocked_id", blockedID))
	return nil
}

func (r *repository) GetBlockedUsers(ctx context.Context, blockerID uint) ([]BlockedUser, error) {
	var blockedUsers []BlockedUser
	if err := r.db.WithContext(ctx).Preload("Blocker").Preload("Blocked").
		Where("blocker_id = ?", blockerID).
		Order("created_at DESC").
		Find(&blockedUsers).Error; err != nil {
		logger.Error("Failed to get blocked users", zap.Error(err))
		return nil, err
	}
	return blockedUsers, nil
}

func (r *repository) CheckUserBlocked(ctx context.Context, blockerID, blockedID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&BlockedUser{}).
		Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check if user is blocked", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func (r *repository) GetBlockedByUsers(ctx context.Context, blockedID uint) ([]BlockedUser, error) {
	var blockedUsers []BlockedUser
	if err := r.db.WithContext(ctx).Preload("Blocker").Preload("Blocked").
		Where("blocked_id = ?", blockedID).
		Order("created_at DESC").
		Find(&blockedUsers).Error; err != nil {
		logger.Error("Failed to get users who blocked this user", zap.Error(err))
		return nil, err
	}
	return blockedUsers, nil
}

// Utility methods

func (r *repository) GetUserByID(ctx context.Context, userID uint) (*user.User, error) {
	var user user.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		logger.Error("Failed to get user by ID", zap.Error(err))
		return nil, err
	}
	return &user, nil
}
