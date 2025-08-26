package friend

import (
	"github.com/gin-gonic/gin"
	"huddle/internal/middleware"
)

// SetupRoutes sets up friend routes
func SetupRoutes(router *gin.RouterGroup, handler *Handler) {
	// Friend routes group (all protected)
	friends := router.Group("/friends")
	friends.Use(middleware.AuthMiddleware())
	{
		// Friend Requests
		friends.POST("/requests", handler.SendFriendRequest)                    // Send friend request
		friends.GET("/requests", handler.GetFriendRequests)                    // Get pending friend requests
		friends.GET("/requests/sent", handler.GetSentFriendRequests)           // Get sent friend requests
		friends.POST("/requests/respond", handler.RespondToFriendRequest)      // Accept/reject friend request
		friends.DELETE("/requests/:request_id", handler.CancelFriendRequest)   // Cancel friend request

		// Friendships
		friends.GET("/", handler.GetFriends)                                   // Get friends list
		friends.DELETE("/:friend_id", handler.RemoveFriend)                    // Remove friend
		friends.GET("/check/:friend_id", handler.CheckFriendship)              // Check if friends

		// Blocked Users
		friends.POST("/block", handler.BlockUser)                              // Block user
		friends.DELETE("/block/:user_id", handler.UnblockUser)                 // Unblock user
		friends.GET("/blocked", handler.GetBlockedUsers)                       // Get blocked users list
		friends.GET("/blocked/check/:user_id", handler.CheckUserBlocked)       // Check if user is blocked
	}
}
