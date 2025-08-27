// User types
export interface User {
  id: number;
  username: string;
  email: string;
  display_name: string;
  bio?: string;
  avatar?: string;
  is_public: boolean;
  last_login?: string;
  created_at: string;
  updated_at: string;
}

// Friend Request types
export interface FriendRequest {
  id: number;
  sender: User;
  receiver: User;
  status: "pending" | "accepted" | "rejected" | "cancelled";
  message?: string;
  created_at: string;
  updated_at: string;
}

export interface FriendRequestList {
  requests: FriendRequest[];
  total: number;
}

// Friendship types
export interface Friendship {
  id: number;
  user: User;
  friend: User;
  created_at: string;
}

export interface FriendList {
  friends: Friendship[];
  total: number;
}

// Blocked User types
export interface BlockedUser {
  id: number;
  blocker: User;
  blocked: User;
  reason?: string;
  created_at: string;
}

export interface BlockedUserList {
  blocked_users: BlockedUser[];
  total: number;
}

// API Response types
export interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
  message?: string;
  timestamp: string;
}

// Search types
export interface UserSearchResult {
  users: User[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}
