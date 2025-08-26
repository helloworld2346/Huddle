package friend

import (
	"context"
	"errors"

	"huddle/pkg/logger"

	"go.uber.org/zap"
)

type service struct {
	repo Repository
}

// NewService creates a new friend service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Friend Requests

func (s *service) SendFriendRequest(ctx context.Context, senderID uint, req *SendFriendRequestRequest) (*FriendRequestResponse, error) {
	// Validate request
	if err := s.ValidateFriendRequest(ctx, senderID, req.ReceiverID); err != nil {
		return nil, err
	}

	// Check if receiver exists
	receiver, err := s.repo.GetUserByID(ctx, req.ReceiverID)
	if err != nil {
		logger.Error("Failed to get receiver user", zap.Error(err))
		return nil, errors.New("receiver not found")
	}

	// Create friend request
	friendRequest, err := s.repo.CreateFriendRequest(ctx, senderID, req.ReceiverID, req.Message)
	if err != nil {
		logger.Error("Failed to create friend request", zap.Error(err))
		return nil, err
	}

	response := &FriendRequestResponse{
		ID:         friendRequest.ID,
		Sender:     friendRequest.Sender,
		Receiver:   friendRequest.Receiver,
		Status:     friendRequest.Status,
		Message:    friendRequest.Message,
		CreatedAt:  friendRequest.CreatedAt,
		UpdatedAt:  friendRequest.UpdatedAt,
	}

	logger.Info("Friend request sent successfully", 
		zap.Uint("sender_id", senderID), 
		zap.Uint("receiver_id", req.ReceiverID),
		zap.String("receiver_username", receiver.Username),
	)

	return response, nil
}

func (s *service) GetFriendRequests(ctx context.Context, userID uint) (*FriendRequestListResponse, error) {
	requests, err := s.repo.GetPendingFriendRequests(ctx, userID)
	if err != nil {
		logger.Error("Failed to get friend requests", zap.Error(err))
		return nil, err
	}

	responses := make([]FriendRequestResponse, len(requests))
	for i, req := range requests {
		responses[i] = FriendRequestResponse{
			ID:         req.ID,
			Sender:     req.Sender,
			Receiver:   req.Receiver,
			Status:     req.Status,
			Message:    req.Message,
			CreatedAt:  req.CreatedAt,
			UpdatedAt:  req.UpdatedAt,
		}
	}

	return &FriendRequestListResponse{
		Requests: responses,
		Total:    len(responses),
	}, nil
}

func (s *service) GetSentFriendRequests(ctx context.Context, userID uint) (*FriendRequestListResponse, error) {
	requests, err := s.repo.GetSentFriendRequests(ctx, userID)
	if err != nil {
		logger.Error("Failed to get sent friend requests", zap.Error(err))
		return nil, err
	}

	responses := make([]FriendRequestResponse, len(requests))
	for i, req := range requests {
		responses[i] = FriendRequestResponse{
			ID:         req.ID,
			Sender:     req.Sender,
			Receiver:   req.Receiver,
			Status:     req.Status,
			Message:    req.Message,
			CreatedAt:  req.CreatedAt,
			UpdatedAt:  req.UpdatedAt,
		}
	}

	return &FriendRequestListResponse{
		Requests: responses,
		Total:    len(responses),
	}, nil
}

func (s *service) RespondToFriendRequest(ctx context.Context, userID uint, req *RespondToFriendRequestRequest) error {
	// Get friend request
	friendRequest, err := s.repo.GetFriendRequestByID(ctx, req.RequestID)
	if err != nil {
		logger.Error("Failed to get friend request", zap.Error(err))
		return errors.New("friend request not found")
	}

	// Check if user is the receiver
	if friendRequest.ReceiverID != userID {
		logger.Error("User is not the receiver of this friend request", 
			zap.Uint("user_id", userID), 
			zap.Uint("receiver_id", friendRequest.ReceiverID),
		)
		return errors.New("unauthorized to respond to this friend request")
	}

	// Check if request is still pending
	if friendRequest.Status != RequestStatusPending {
		logger.Error("Friend request is not pending", 
			zap.Uint("request_id", req.RequestID), 
			zap.String("status", friendRequest.Status),
		)
		return errors.New("friend request is not pending")
	}

	// Update status
	var newStatus string
	switch req.Action {
	case "accept":
		newStatus = RequestStatusAccepted
		// Create friendship
		if _, err := s.repo.CreateFriendship(ctx, friendRequest.SenderID, friendRequest.ReceiverID); err != nil {
			logger.Error("Failed to create friendship", zap.Error(err))
			return errors.New("failed to create friendship")
		}
	case "reject":
		newStatus = RequestStatusRejected
	default:
		return errors.New("invalid action")
	}

	if err := s.repo.UpdateFriendRequestStatus(ctx, req.RequestID, newStatus); err != nil {
		logger.Error("Failed to update friend request status", zap.Error(err))
		return err
	}

	logger.Info("Friend request responded", 
		zap.Uint("request_id", req.RequestID), 
		zap.String("action", req.Action),
		zap.Uint("sender_id", friendRequest.SenderID),
		zap.Uint("receiver_id", friendRequest.ReceiverID),
	)

	return nil
}

func (s *service) CancelFriendRequest(ctx context.Context, userID uint, requestID uint) error {
	// Get friend request
	friendRequest, err := s.repo.GetFriendRequestByID(ctx, requestID)
	if err != nil {
		logger.Error("Failed to get friend request", zap.Error(err))
		return errors.New("friend request not found")
	}

	// Check if user is the sender
	if friendRequest.SenderID != userID {
		logger.Error("User is not the sender of this friend request", 
			zap.Uint("user_id", userID), 
			zap.Uint("sender_id", friendRequest.SenderID),
		)
		return errors.New("unauthorized to cancel this friend request")
	}

	// Check if request is still pending
	if friendRequest.Status != RequestStatusPending {
		logger.Error("Friend request is not pending", 
			zap.Uint("request_id", requestID), 
			zap.String("status", friendRequest.Status),
		)
		return errors.New("friend request is not pending")
	}

	if err := s.repo.UpdateFriendRequestStatus(ctx, requestID, RequestStatusCancelled); err != nil {
		logger.Error("Failed to cancel friend request", zap.Error(err))
		return err
	}

	logger.Info("Friend request cancelled", 
		zap.Uint("request_id", requestID),
		zap.Uint("sender_id", friendRequest.SenderID),
		zap.Uint("receiver_id", friendRequest.ReceiverID),
	)

	return nil
}

// Friendships

func (s *service) GetFriends(ctx context.Context, userID uint) (*FriendListResponse, error) {
	friendships, err := s.repo.GetUserFriends(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user friends", zap.Error(err))
		return nil, err
	}

	responses := make([]FriendshipResponse, len(friendships))
	for i, friendship := range friendships {
		responses[i] = FriendshipResponse{
			ID:        friendship.ID,
			User:      friendship.User,
			Friend:    friendship.Friend,
			CreatedAt: friendship.CreatedAt,
		}
	}

	return &FriendListResponse{
		Friends: responses,
		Total:   len(responses),
	}, nil
}

func (s *service) RemoveFriend(ctx context.Context, userID uint, friendID uint) error {
	// Validate friendship action
	if err := s.ValidateFriendshipAction(ctx, userID, friendID); err != nil {
		return err
	}

	// Check if friendship exists
	exists, err := s.repo.CheckFriendshipExists(ctx, userID, friendID)
	if err != nil {
		logger.Error("Failed to check friendship exists", zap.Error(err))
		return err
	}

	if !exists {
		logger.Error("Friendship does not exist", 
			zap.Uint("user_id", userID), 
			zap.Uint("friend_id", friendID),
		)
		return errors.New("friendship not found")
	}

	if err := s.repo.DeleteFriendship(ctx, userID, friendID); err != nil {
		logger.Error("Failed to remove friend", zap.Error(err))
		return err
	}

	logger.Info("Friend removed", 
		zap.Uint("user_id", userID), 
		zap.Uint("friend_id", friendID),
	)

	return nil
}

func (s *service) CheckFriendship(ctx context.Context, userID, friendID uint) (bool, error) {
	return s.repo.CheckFriendshipExists(ctx, userID, friendID)
}

// Blocked Users

func (s *service) BlockUser(ctx context.Context, blockerID uint, req *BlockUserRequest) (*BlockedUserResponse, error) {
	// Check if user exists
	blockedUser, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		logger.Error("Failed to get user to block", zap.Error(err))
		return nil, errors.New("user not found")
	}

	// Check if already blocked
	alreadyBlocked, err := s.repo.CheckUserBlocked(ctx, blockerID, req.UserID)
	if err != nil {
		logger.Error("Failed to check if user is already blocked", zap.Error(err))
		return nil, err
	}

	if alreadyBlocked {
		logger.Error("User is already blocked", 
			zap.Uint("blocker_id", blockerID), 
			zap.Uint("blocked_id", req.UserID),
		)
		return nil, errors.New("user is already blocked")
	}

	// Block user
	blocked, err := s.repo.BlockUser(ctx, blockerID, req.UserID, req.Reason)
	if err != nil {
		logger.Error("Failed to block user", zap.Error(err))
		return nil, err
	}

	response := &BlockedUserResponse{
		ID:         blocked.ID,
		Blocker:    blocked.Blocker,
		Blocked:    blocked.Blocked,
		Reason:     blocked.Reason,
		CreatedAt:  blocked.CreatedAt,
	}

	logger.Info("User blocked successfully", 
		zap.Uint("blocker_id", blockerID), 
		zap.Uint("blocked_id", req.UserID),
		zap.String("blocked_username", blockedUser.Username),
	)

	return response, nil
}

func (s *service) UnblockUser(ctx context.Context, blockerID uint, blockedID uint) error {
	// Check if user is blocked
	blocked, err := s.repo.CheckUserBlocked(ctx, blockerID, blockedID)
	if err != nil {
		logger.Error("Failed to check if user is blocked", zap.Error(err))
		return err
	}

	if !blocked {
		logger.Error("User is not blocked", 
			zap.Uint("blocker_id", blockerID), 
			zap.Uint("blocked_id", blockedID),
		)
		return errors.New("user is not blocked")
	}

	if err := s.repo.UnblockUser(ctx, blockerID, blockedID); err != nil {
		logger.Error("Failed to unblock user", zap.Error(err))
		return err
	}

	logger.Info("User unblocked successfully", 
		zap.Uint("blocker_id", blockerID), 
		zap.Uint("blocked_id", blockedID),
	)

	return nil
}

func (s *service) GetBlockedUsers(ctx context.Context, blockerID uint) (*BlockedUserListResponse, error) {
	blockedUsers, err := s.repo.GetBlockedUsers(ctx, blockerID)
	if err != nil {
		logger.Error("Failed to get blocked users", zap.Error(err))
		return nil, err
	}

	responses := make([]BlockedUserResponse, len(blockedUsers))
	for i, blocked := range blockedUsers {
		responses[i] = BlockedUserResponse{
			ID:         blocked.ID,
			Blocker:    blocked.Blocker,
			Blocked:    blocked.Blocked,
			Reason:     blocked.Reason,
			CreatedAt:  blocked.CreatedAt,
		}
	}

	return &BlockedUserListResponse{
		BlockedUsers: responses,
		Total:        len(responses),
	}, nil
}

func (s *service) CheckUserBlocked(ctx context.Context, blockerID, blockedID uint) (bool, error) {
	return s.repo.CheckUserBlocked(ctx, blockerID, blockedID)
}

// Validation methods

func (s *service) ValidateFriendRequest(ctx context.Context, senderID, receiverID uint) error {
	// Check if sender and receiver are the same
	if senderID == receiverID {
		return errors.New("cannot send friend request to yourself")
	}

	// Check if receiver exists
	_, err := s.repo.GetUserByID(ctx, receiverID)
	if err != nil {
		return errors.New("receiver not found")
	}

	// Check if already friends
	alreadyFriends, err := s.repo.CheckFriendshipExists(ctx, senderID, receiverID)
	if err != nil {
		logger.Error("Failed to check friendship exists", zap.Error(err))
		return err
	}

	if alreadyFriends {
		return errors.New("already friends")
	}

	// Check if there's already a pending request
	existingRequest, err := s.repo.GetFriendRequestByUsers(ctx, senderID, receiverID)
	if err == nil && existingRequest != nil {
		if existingRequest.Status == RequestStatusPending {
			return errors.New("friend request already pending")
		}
	}

	// Check if blocked
	blocked, err := s.repo.CheckUserBlocked(ctx, senderID, receiverID)
	if err != nil {
		logger.Error("Failed to check if user is blocked", zap.Error(err))
		return err
	}

	if blocked {
		return errors.New("cannot send friend request to blocked user")
	}

	// Check if blocked by receiver
	blockedByReceiver, err := s.repo.CheckUserBlocked(ctx, receiverID, senderID)
	if err != nil {
		logger.Error("Failed to check if blocked by receiver", zap.Error(err))
		return err
	}

	if blockedByReceiver {
		return errors.New("cannot send friend request to user who blocked you")
	}

	return nil
}

func (s *service) ValidateFriendshipAction(ctx context.Context, userID, targetID uint) error {
	// Check if target user exists
	_, err := s.repo.GetUserByID(ctx, targetID)
	if err != nil {
		return errors.New("target user not found")
	}

	// Check if users are the same
	if userID == targetID {
		return errors.New("cannot perform action on yourself")
	}

	return nil
}
