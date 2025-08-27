const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

export interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  message?: string;
  error?: {
    code: string;
    message: string;
  };
  timestamp?: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  display_name: string;
  password: string;
  bio?: string;
  is_public?: boolean;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface AuthResponse {
  user: {
    id: number;
    username: string;
    email: string;
    display_name: string;
    bio: string;
    avatar: string;
    is_public: boolean;
    last_login: string | null;
    created_at: string;
    updated_at: string;
  };
  tokens: {
    access_token: string;
    refresh_token: string;
    expires_in: number;
  };
  message: string;
}

// Generic API call function
async function apiCall<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const url = `${API_BASE_URL}${endpoint}`;

  const defaultOptions: RequestInit = {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  };

  const response = await fetch(url, {
    ...defaultOptions,
    ...options,
  });

  const data = await response.json();
  return data;
}

// Auth API functions
export const authApi = {
  register: async (
    data: RegisterRequest
  ): Promise<ApiResponse<AuthResponse>> => {
    return apiCall<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  login: async (data: LoginRequest): Promise<ApiResponse<AuthResponse>> => {
    return apiCall<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  forgotPassword: async (email: string): Promise<ApiResponse> => {
    return apiCall("/auth/forgot-password", {
      method: "POST",
      body: JSON.stringify({ email }),
    });
  },

  resetPassword: async (
    token: string,
    newPassword: string
  ): Promise<ApiResponse> => {
    return apiCall("/auth/reset-password", {
      method: "POST",
      body: JSON.stringify({ token, new_password: newPassword }),
    });
  },

  refreshToken: async (
    refreshToken: string
  ): Promise<ApiResponse<AuthResponse>> => {
    return apiCall<AuthResponse>("/auth/refresh", {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  },

  logout: async (
    refreshToken: string,
    accessToken?: string
  ): Promise<ApiResponse> => {
    const headers: Record<string, string> = {};
    if (accessToken) {
      headers["Authorization"] = `Bearer ${accessToken}`;
    }

    return apiCall("/auth/logout", {
      method: "POST",
      headers,
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  },
};

// Dashboard API interfaces
export interface ConversationResponse {
  id: number;
  name: string;
  type: "direct" | "group";
  last_message?: {
    content: string;
    sender: string;
    timestamp: string;
  };
  unread_count: number;
  participants: UserResponse[];
  updated_at: string;
}

export interface MessageResponse {
  id: number;
  content: string;
  sender: UserResponse;
  timestamp: string;
  type: "text" | "file" | "image";
  file_url?: string;
  file_name?: string;
}

export interface UserResponse {
  id: number;
  username: string;
  display_name: string;
  avatar: string;
  is_online: boolean;
  last_seen?: string;
}

// Dashboard API functions
export const dashboardApi = {
  // Get user conversations
  getConversations: async (
    accessToken: string
  ): Promise<ApiResponse<ConversationResponse[]>> => {
    return apiCall<ConversationResponse[]>("/conversations", {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
  },

  // Get conversation messages
  getMessages: async (
    conversationId: number,
    accessToken: string
  ): Promise<ApiResponse<MessageResponse[]>> => {
    return apiCall<MessageResponse[]>(
      `/conversations/${conversationId}/messages`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    );
  },

  // Send message
  sendMessage: async (
    conversationId: number,
    content: string,
    accessToken: string
  ): Promise<ApiResponse<MessageResponse>> => {
    return apiCall<MessageResponse>(
      `/conversations/${conversationId}/messages`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({ content }),
      }
    );
  },

  // Create new conversation
  createConversation: async (
    data: {
      name?: string;
      type: "direct" | "group";
      participant_ids: number[];
    },
    accessToken: string
  ): Promise<ApiResponse<ConversationResponse>> => {
    return apiCall<ConversationResponse>("/conversations", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
      body: JSON.stringify(data),
    });
  },

  // Search users
  searchUsers: async (
    query: string,
    accessToken: string
  ): Promise<ApiResponse<UserResponse[]>> => {
    return apiCall<UserResponse[]>(
      `/users/search?q=${encodeURIComponent(query)}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    );
  },
};
