package friend

import (
	"context"
	"huddle/internal/user"
)

type Repository interface {
	CreateFriendRequest(ctx context.Context, senderID, receiverID uint, message string) (*FriendRequest, error)
	GetFriendRequestByID(ctx context.Context, requestID uint) (*FriendRequest, error)
	GetFriendRequestByUsers(ctx context.Context, senderID, receiverID uint) (*FriendRequest, error)
	GetPendingFriendRequests(ctx context.Context, receiverID uint) ([]FriendRequest, error)
	GetSentFriendRequests(ctx context.Context, senderID uint) ([]FriendRequest, error)
	UpdateFriendRequestStatus(ctx context.Context, requestID uint, status string) error
	DeleteFriendRequest(ctx context.Context, requestID uint) error
	CreateFriendship(ctx context.Context, userID, friendID uint) (*Friendship, error)
	GetFriendship(ctx context.Context, userID, friendID uint) (*Friendship, error)
	GetUserFriends(ctx context.Context, userID uint) ([]Friendship, error)
	DeleteFriendship(ctx context.Context, userID, friendID uint) error
	CheckFriendshipExists(ctx context.Context, userID, friendID uint) (bool, error)
	BlockUser(ctx context.Context, blockerID, blockedID uint, reason string) (*BlockedUser, error)
	UnblockUser(ctx context.Context, blockerID, blockedID uint) error
	GetBlockedUsers(ctx context.Context, blockerID uint) ([]BlockedUser, error)
	CheckUserBlocked(ctx context.Context, blockerID, blockedID uint) (bool, error)
	GetBlockedByUsers(ctx context.Context, blockedID uint) ([]BlockedUser, error)
	GetUserByID(ctx context.Context, userID uint) (*user.User, error)
}

type Service interface {
	SendFriendRequest(ctx context.Context, senderID uint, req *SendFriendRequestRequest) (*FriendRequestResponse, error)
	GetFriendRequests(ctx context.Context, userID uint) (*FriendRequestListResponse, error)
	GetSentFriendRequests(ctx context.Context, userID uint) (*FriendRequestListResponse, error)
	RespondToFriendRequest(ctx context.Context, userID uint, req *RespondToFriendRequestRequest) error
	CancelFriendRequest(ctx context.Context, userID uint, requestID uint) error
	GetFriends(ctx context.Context, userID uint) (*FriendListResponse, error)
	RemoveFriend(ctx context.Context, userID uint, friendID uint) error
	CheckFriendship(ctx context.Context, userID, friendID uint) (bool, error)
	BlockUser(ctx context.Context, blockerID uint, req *BlockUserRequest) (*BlockedUserResponse, error)
	UnblockUser(ctx context.Context, blockerID uint, blockedID uint) error
	GetBlockedUsers(ctx context.Context, blockerID uint) (*BlockedUserListResponse, error)
	CheckUserBlocked(ctx context.Context, blockerID, blockedID uint) (bool, error)
	ValidateFriendRequest(ctx context.Context, senderID, receiverID uint) error
	ValidateFriendshipAction(ctx context.Context, userID, targetID uint) error
}