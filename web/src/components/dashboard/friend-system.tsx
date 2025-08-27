"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Users,
  UserPlus,
  UserCheck,
  UserX,
  Search,
  Send,
  Check,
  X,
  Loader2,
  User as UserIcon,
} from "lucide-react";
import { toast } from "sonner";
import { friendAPI, userAPI } from "@/lib/api";
import type { User, FriendRequest, Friendship } from "@/lib/types";

interface FriendSystemProps {
  currentUser: User;
}

export function FriendSystem({ currentUser }: FriendSystemProps) {
  const [activeTab, setActiveTab] = useState<"friends" | "requests" | "search">(
    "friends"
  );
  const [isLoading, setIsLoading] = useState(false);

  // Friends state
  const [friends, setFriends] = useState<Friendship[]>([]);
  const [friendsLoading, setFriendsLoading] = useState(false);

  // Friend requests state
  const [friendRequests, setFriendRequests] = useState<FriendRequest[]>([]);
  const [sentRequests, setSentRequests] = useState<FriendRequest[]>([]);
  const [requestsLoading, setRequestsLoading] = useState(false);

  // Search state
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<User[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [requestMessage, setRequestMessage] = useState("");

  // Helper function to check if user already has a friend request with current user
  const hasFriendRequestWithUser = (userId: number) => {
    return sentRequests.some((request) => request.receiver.id === userId);
  };

  // Helper function to check if user has received a friend request from current user
  const hasReceivedRequestFromUser = (userId: number) => {
    return friendRequests.some((request) => request.sender.id === userId);
  };

  // Helper function to check if user is already a friend
  const isAlreadyFriend = (userId: number) => {
    return friends.some((friendship) => friendship.friend.id === userId);
  };

  // Load friends
  const loadFriends = async () => {
    setFriendsLoading(true);
    try {
      const response = await friendAPI.getFriends();
      if (response.success && response.data) {
        setFriends(response.data.friends);
      }
    } catch {
      toast.error("Failed to load friends");
    } finally {
      setFriendsLoading(false);
    }
  };

  // Load friend requests
  const loadFriendRequests = async () => {
    setRequestsLoading(true);
    try {
      const [receivedResponse, sentResponse] = await Promise.all([
        friendAPI.getFriendRequests(),
        friendAPI.getSentFriendRequests(),
      ]);

      if (receivedResponse.success && receivedResponse.data) {
        setFriendRequests(receivedResponse.data.requests);
      }
      if (sentResponse.success && sentResponse.data) {
        setSentRequests(sentResponse.data.requests);
      }
    } catch {
      toast.error("Failed to load friend requests");
    } finally {
      setRequestsLoading(false);
    }
  };

  // Search users
  const searchUsers = async (query: string) => {
    if (!query.trim()) {
      setSearchResults([]);
      setHasSearched(false);
      return;
    }

    setSearchLoading(true);
    setHasSearched(true);
    try {
      const response = await userAPI.searchUsers(query);
      if (response.success && response.data) {
        setSearchResults(
          response.data.users.filter((user: User) => user.id !== currentUser.id)
        );
      }
    } catch {
      toast.error("Failed to search users");
    } finally {
      setSearchLoading(false);
    }
  };

  // Debounced search effect
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      searchUsers(searchQuery);
    }, 500); // 500ms delay for better UX

    return () => clearTimeout(timeoutId);
  }, [searchQuery]);

  // Send friend request
  const sendFriendRequest = async () => {
    if (!selectedUser) return;

    setIsLoading(true);
    try {
      const response = await friendAPI.sendFriendRequest({
        receiver_id: selectedUser.id,
        message: requestMessage.trim() || undefined,
      });

      if (response.success) {
        toast.success("Friend request sent!");
        setSelectedUser(null);
        setRequestMessage("");
        setSearchQuery("");
        setSearchResults([]);
        setHasSearched(false);
        // Refresh both friends and requests to update status
        loadFriends();
        loadFriendRequests();
      } else {
        // Handle specific error cases
        if (response.error?.message?.includes("duplicate key")) {
          toast.error("Friend request already sent to this user");
        } else {
          throw new Error(response.error?.message || "Failed to send request");
        }
      }
    } catch {
      toast.error("Failed to send friend request");
    } finally {
      setIsLoading(false);
    }
  };

  // Respond to friend request
  const respondToRequest = async (
    requestId: number,
    action: "accept" | "reject"
  ) => {
    setIsLoading(true);
    try {
      const response = await friendAPI.respondToFriendRequest({
        request_id: requestId,
        action,
      });

      if (response.success) {
        toast.success(`Friend request ${action}ed!`);
        loadFriendRequests(); // Refresh requests
        if (action === "accept") {
          loadFriends(); // Refresh friends list
        }
      } else {
        throw new Error(
          response.error?.message || "Failed to respond to request"
        );
      }
    } catch {
      toast.error("Failed to respond to friend request");
    } finally {
      setIsLoading(false);
    }
  };

  // Remove friend
  const removeFriend = async (friendId: number) => {
    setIsLoading(true);
    try {
      const response = await friendAPI.removeFriend(friendId);
      if (response.success) {
        toast.success("Friend removed");
        loadFriends(); // Refresh friends list
      } else {
        throw new Error(response.error?.message || "Failed to remove friend");
      }
    } catch {
      toast.error("Failed to remove friend");
    } finally {
      setIsLoading(false);
    }
  };

  // Load data when tab changes
  useEffect(() => {
    if (activeTab === "friends") {
      loadFriends();
    } else if (activeTab === "requests") {
      loadFriendRequests();
    } else if (activeTab === "search") {
      // Load friend requests when searching to check status
      loadFriendRequests();
    }
  }, [activeTab]);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <Users className="h-5 w-5" />
          <span>Friends</span>
        </CardTitle>
        <CardDescription>
          Manage your friends and friend requests
        </CardDescription>
      </CardHeader>
      <CardContent>
        {/* Tabs */}
        <div className="flex space-x-1 mb-6">
          <Button
            variant={activeTab === "friends" ? "default" : "outline"}
            size="sm"
            onClick={() => setActiveTab("friends")}
          >
            Friends ({friends.length})
          </Button>
          <Button
            variant={activeTab === "requests" ? "default" : "outline"}
            size="sm"
            onClick={() => setActiveTab("requests")}
          >
            Requests ({friendRequests.length + sentRequests.length})
          </Button>
          <Button
            variant={activeTab === "search" ? "default" : "outline"}
            size="sm"
            onClick={() => setActiveTab("search")}
          >
            <Search className="h-4 w-4 mr-2" />
            Add Friends
          </Button>
        </div>

        {/* Friends Tab */}
        {activeTab === "friends" && (
          <div className="space-y-4">
            {friendsLoading ? (
              <div className="flex items-center justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin" />
              </div>
            ) : friends.length === 0 ? (
              <div className="text-center py-12 text-muted-foreground">
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Users className="h-8 w-8 text-primary" />
                </div>
                <h3 className="text-lg font-medium mb-2">No friends yet</h3>
                <p className="text-sm mb-4">
                  Start connecting with other users
                </p>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setActiveTab("search")}
                >
                  <Search className="h-4 w-4 mr-2" />
                  Find Friends
                </Button>
              </div>
            ) : (
              <>
                <div className="flex items-center justify-between">
                  <h3 className="font-semibold text-lg">
                    Your Friends ({friends.length})
                  </h3>
                </div>
                <div className="grid gap-3">
                  {friends.map((friendship) => (
                    <div
                      key={friendship.id}
                      className="group flex items-center justify-between p-4 border rounded-xl hover:border-primary/50 hover:bg-primary/5 transition-all duration-200"
                    >
                      <div className="flex items-center space-x-4">
                        <div className="w-12 h-12 bg-gradient-to-br from-green-500/20 to-green-500/10 rounded-full flex items-center justify-center">
                          <UserCheck className="h-6 w-6 text-green-600" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="font-semibold text-base truncate">
                            {friendship.friend.display_name}
                          </p>
                          <p className="text-sm text-muted-foreground">
                            @{friendship.friend.username}
                          </p>
                          {friendship.friend.bio && (
                            <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                              {friendship.friend.bio}
                            </p>
                          )}
                        </div>
                      </div>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => removeFriend(friendship.friend.id)}
                        disabled={isLoading}
                        className="opacity-0 group-hover:opacity-100 transition-opacity duration-200"
                      >
                        <UserX className="h-4 w-4 mr-2" />
                        Remove
                      </Button>
                    </div>
                  ))}
                </div>
              </>
            )}
          </div>
        )}

        {/* Requests Tab */}
        {activeTab === "requests" && (
          <div className="space-y-6">
            {requestsLoading ? (
              <div className="flex items-center justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin" />
              </div>
            ) : friendRequests.length === 0 && sentRequests.length === 0 ? (
              <div className="text-center py-12 text-muted-foreground">
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4">
                  <UserPlus className="h-8 w-8 text-primary" />
                </div>
                <h3 className="text-lg font-medium mb-2">No friend requests</h3>
                <p className="text-sm mb-4">
                  You have no pending or sent requests
                </p>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setActiveTab("search")}
                >
                  <Search className="h-4 w-4 mr-2" />
                  Find Friends
                </Button>
              </div>
            ) : (
              <>
                {/* Received Requests */}
                <div>
                  <h3 className="font-medium mb-3">Received Requests</h3>
                  {friendRequests.length === 0 ? (
                    <p className="text-muted-foreground text-sm">
                      No pending requests
                    </p>
                  ) : (
                    <div className="space-y-3">
                      {friendRequests.map((request) => (
                        <div
                          key={request.id}
                          className="flex items-center justify-between p-4 border rounded-lg"
                        >
                          <div className="flex items-center space-x-3">
                            <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                              <UserPlus className="h-5 w-5 text-primary" />
                            </div>
                            <div>
                              <p className="font-medium">
                                {request.sender.display_name}
                              </p>
                              <p className="text-sm text-muted-foreground">
                                @{request.sender.username}
                              </p>
                              {request.message && (
                                <p className="text-sm text-muted-foreground mt-1">
                                  {request.message}
                                </p>
                              )}
                            </div>
                          </div>
                          <div className="flex space-x-2">
                            <Button
                              size="sm"
                              onClick={() =>
                                respondToRequest(request.id, "accept")
                              }
                              disabled={isLoading}
                            >
                              <Check className="h-4 w-4 mr-2" />
                              Accept
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() =>
                                respondToRequest(request.id, "reject")
                              }
                              disabled={isLoading}
                            >
                              <X className="h-4 w-4 mr-2" />
                              Reject
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                {/* Sent Requests */}
                <div>
                  <h3 className="font-medium mb-3">Sent Requests</h3>
                  {sentRequests.length === 0 ? (
                    <p className="text-muted-foreground text-sm">
                      No sent requests
                    </p>
                  ) : (
                    <div className="space-y-3">
                      {sentRequests.map((request) => (
                        <div
                          key={request.id}
                          className="flex items-center justify-between p-4 border rounded-lg"
                        >
                          <div className="flex items-center space-x-3">
                            <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                              <UserPlus className="h-5 w-5 text-primary" />
                            </div>
                            <div>
                              <p className="font-medium">
                                {request.receiver.display_name}
                              </p>
                              <p className="text-sm text-muted-foreground">
                                @{request.receiver.username}
                              </p>
                              <p className="text-sm text-muted-foreground">
                                Status: {request.status}
                              </p>
                            </div>
                          </div>
                          {request.status === "pending" && (
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={async () => {
                                try {
                                  const response =
                                    await friendAPI.cancelFriendRequest(
                                      request.id
                                    );
                                  if (response.success) {
                                    toast.success("Friend request cancelled");
                                    loadFriendRequests(); // Refresh the list
                                  } else {
                                    throw new Error(
                                      response.error?.message ||
                                        "Failed to cancel request"
                                    );
                                  }
                                } catch {
                                  toast.error(
                                    "Failed to cancel friend request"
                                  );
                                }
                              }}
                              disabled={isLoading}
                              className="opacity-0 group-hover:opacity-100 transition-opacity duration-200"
                            >
                              Cancel
                            </Button>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </>
            )}
          </div>
        )}

        {/* Search Tab */}
        {activeTab === "search" && (
          <div className="space-y-4">
            <div className="flex space-x-2">
              <div className="flex-1 relative">
                <Input
                  placeholder="Search users by username or display name..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className={searchLoading ? "pr-10" : ""}
                />
                {searchLoading && (
                  <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                    <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
                  </div>
                )}
              </div>
            </div>

            {/* Search Results */}
            {!hasSearched && searchQuery.trim() === "" ? (
              <div className="text-center py-12 text-muted-foreground">
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Search className="h-8 w-8 text-primary" />
                </div>
                <h3 className="text-lg font-medium mb-2">Find Friends</h3>
                <p className="text-sm">Start typing to search for users</p>
              </div>
            ) : (
              <div className="space-y-4">
                {searchResults.length > 0 ? (
                  <>
                    <div className="flex items-center justify-between">
                      <h3 className="font-semibold text-lg">
                        Found {searchResults.length} user
                        {searchResults.length !== 1 ? "s" : ""}
                      </h3>
                      <div className="text-sm text-muted-foreground">
                        Click &quot;Send Request&quot; to add as friend
                      </div>
                    </div>
                    <div className="grid gap-3 max-h-96 overflow-y-auto">
                      {searchResults.map((user) => (
                        <div
                          key={user.id}
                          className="group flex items-center justify-between p-4 border rounded-xl hover:border-primary/50 hover:bg-primary/5 transition-all duration-200"
                        >
                          <div className="flex items-center space-x-4">
                            <div className="w-12 h-12 bg-gradient-to-br from-primary/20 to-primary/10 rounded-full flex items-center justify-center">
                              <UserIcon className="h-6 w-6 text-primary" />
                            </div>
                            <div className="flex-1 min-w-0">
                              <div className="flex items-center space-x-2">
                                <p className="font-semibold text-base truncate">
                                  {user.display_name}
                                </p>
                                {user.is_public && (
                                  <span className="px-2 py-1 text-xs bg-green-100 text-green-700 rounded-full">
                                    Public
                                  </span>
                                )}
                              </div>
                              <p className="text-sm text-muted-foreground">
                                @{user.username}
                              </p>
                              {user.bio && (
                                <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                                  {user.bio}
                                </p>
                              )}
                            </div>
                          </div>
                          {isAlreadyFriend(user.id) ? (
                            <span className="px-3 py-1 text-xs bg-green-100 text-green-700 rounded-full">
                              Already Friends
                            </span>
                          ) : hasReceivedRequestFromUser(user.id) ? (
                            <div className="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                              <span className="px-3 py-1 text-xs bg-blue-100 text-blue-700 rounded-full">
                                Request Received
                              </span>
                              <Button
                                size="sm"
                                onClick={async () => {
                                  try {
                                    const request = friendRequests.find(
                                      (req) => req.sender.id === user.id
                                    );
                                    if (request) {
                                      const response =
                                        await friendAPI.respondToFriendRequest({
                                          request_id: request.id,
                                          action: "accept",
                                        });
                                      if (response.success) {
                                        toast.success(
                                          "Friend request accepted!"
                                        );
                                        loadFriends();
                                        loadFriendRequests();
                                      } else {
                                        throw new Error(
                                          response.error?.message ||
                                            "Failed to accept request"
                                        );
                                      }
                                    }
                                  } catch (error) {
                                    toast.error(
                                      "Failed to accept friend request"
                                    );
                                  }
                                }}
                                disabled={isLoading}
                              >
                                <Check className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={async () => {
                                  try {
                                    const request = friendRequests.find(
                                      (req) => req.sender.id === user.id
                                    );
                                    if (request) {
                                      const response =
                                        await friendAPI.respondToFriendRequest({
                                          request_id: request.id,
                                          action: "reject",
                                        });
                                      if (response.success) {
                                        toast.success(
                                          "Friend request rejected"
                                        );
                                        loadFriendRequests();
                                      } else {
                                        throw new Error(
                                          response.error?.message ||
                                            "Failed to reject request"
                                        );
                                      }
                                    }
                                  } catch (error) {
                                    toast.error(
                                      "Failed to reject friend request"
                                    );
                                  }
                                }}
                                disabled={isLoading}
                              >
                                <X className="h-4 w-4" />
                              </Button>
                            </div>
                          ) : hasFriendRequestWithUser(user.id) ? (
                            <div className="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                              <span className="px-3 py-1 text-xs bg-yellow-100 text-yellow-700 rounded-full">
                                Request Sent
                              </span>
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={async () => {
                                  try {
                                    const request = sentRequests.find(
                                      (req) => req.receiver.id === user.id
                                    );
                                    if (request) {
                                      const response =
                                        await friendAPI.cancelFriendRequest(
                                          request.id
                                        );
                                      if (response.success) {
                                        toast.success(
                                          "Friend request cancelled"
                                        );
                                        loadFriends();
                                        loadFriendRequests();
                                      } else {
                                        throw new Error(
                                          response.error?.message ||
                                            "Failed to cancel request"
                                        );
                                      }
                                    }
                                  } catch (error) {
                                    toast.error(
                                      "Failed to cancel friend request"
                                    );
                                  }
                                }}
                                disabled={isLoading}
                              >
                                <X className="h-4 w-4" />
                              </Button>
                            </div>
                          ) : (
                            <Button
                              size="sm"
                              onClick={() => setSelectedUser(user)}
                              disabled={isLoading}
                              className="opacity-0 group-hover:opacity-100 transition-opacity duration-200"
                            >
                              <Send className="h-4 w-4 mr-2" />
                              Send Request
                            </Button>
                          )}
                        </div>
                      ))}
                    </div>
                  </>
                ) : searchQuery.trim() !== "" ? (
                  <div className="text-center py-12 text-muted-foreground">
                    <div className="w-16 h-16 bg-orange-100 dark:bg-orange-900/20 rounded-full flex items-center justify-center mx-auto mb-4">
                      <Search className="h-8 w-8 text-orange-500" />
                    </div>
                    <h3 className="text-lg font-medium mb-2">No users found</h3>
                    <p className="text-sm mb-4">
                      No users match &quot;{searchQuery}&quot;
                    </p>
                    <div className="text-xs space-y-1">
                      <p>• Try searching with a different username</p>
                      <p>• Check if the user has a public profile</p>
                    </div>
                  </div>
                ) : null}
              </div>
            )}

            {/* Send Request Modal */}
            {selectedUser && (
              <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                <div className="bg-background p-6 rounded-lg w-full max-w-md mx-4">
                  <h3 className="font-medium mb-4">
                    Send friend request to {selectedUser.display_name}
                  </h3>
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="message">Message (optional)</Label>
                      <Textarea
                        id="message"
                        value={requestMessage}
                        onChange={(e) => setRequestMessage(e.target.value)}
                        placeholder="Add a personal message..."
                        rows={3}
                      />
                    </div>
                    <div className="flex space-x-2">
                      <Button
                        variant="outline"
                        onClick={() => setSelectedUser(null)}
                        disabled={isLoading}
                      >
                        Cancel
                      </Button>
                      <Button onClick={sendFriendRequest} disabled={isLoading}>
                        {isLoading ? (
                          <Loader2 className="h-4 w-4 animate-spin mr-2" />
                        ) : (
                          <Send className="h-4 w-4 mr-2" />
                        )}
                        Send Request
                      </Button>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
