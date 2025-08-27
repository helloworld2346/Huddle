import axios from "axios";

// API base URL - có thể config từ environment
const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

// Create axios instance
export const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor để thêm auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("access_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor để handle errors
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Handle 401 errors - token expired
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const refreshToken = localStorage.getItem("refresh_token");
      if (refreshToken) {
        try {
          const response = await api.post("/auth/refresh", {
            refresh_token: refreshToken,
          });

          if (response.data.success && response.data.data) {
            const { tokens } = response.data.data;
            localStorage.setItem("access_token", tokens.access_token);
            originalRequest.headers.Authorization = `Bearer ${tokens.access_token}`;
          } else {
            throw new Error("Failed to refresh token");
          }
          return api(originalRequest);
        } catch (refreshError) {
          // Refresh token expired, redirect to login
          localStorage.removeItem("access_token");
          localStorage.removeItem("refresh_token");
          window.location.href = "/auth/login";
          return Promise.reject(refreshError);
        }
      }
    }

    return Promise.reject(error);
  }
);

// Auth API functions
export const authAPI = {
  login: async (credentials: { username: string; password: string }) => {
    const response = await api.post("/auth/login", credentials);
    return response.data;
  },

  register: async (userData: {
    username: string;
    email: string;
    display_name: string;
    password: string;
  }) => {
    const response = await api.post("/auth/register", userData);
    return response.data;
  },

  logout: async () => {
    const refreshToken = localStorage.getItem("refresh_token");
    if (!refreshToken) {
      throw new Error("No refresh token found");
    }
    const response = await api.post("/auth/logout", {
      refresh_token: refreshToken,
    });
    return response.data;
  },

  refreshToken: async (refreshToken: string) => {
    const response = await api.post("/auth/refresh", {
      refresh_token: refreshToken,
    });
    return response.data;
  },
};

// User API functions
export const userAPI = {
  getCurrentUser: async () => {
    const response = await api.get("/users/me");
    return response.data;
  },

  updateProfile: async (userData: {
    display_name?: string;
    bio?: string;
    avatar?: string;
    is_public?: boolean;
  }) => {
    const response = await api.put("/users/me", userData);
    return response.data;
  },

  searchUsers: async (query: string) => {
    const response = await api.get(
      `/users/search?q=${encodeURIComponent(query)}`
    );
    return response.data;
  },

  getUserByUsername: async (username: string) => {
    const response = await api.get(`/users/username/${username}`);
    return response.data;
  },
};

// Friend API functions
export const friendAPI = {
  // Friend Requests
  sendFriendRequest: async (data: {
    receiver_id: number;
    message?: string;
  }) => {
    const response = await api.post("/friends/requests", data);
    return response.data;
  },

  getFriendRequests: async () => {
    const response = await api.get("/friends/requests");
    return response.data;
  },

  getSentFriendRequests: async () => {
    const response = await api.get("/friends/requests/sent");
    return response.data;
  },

  respondToFriendRequest: async (data: {
    request_id: number;
    action: "accept" | "reject";
  }) => {
    const response = await api.post("/friends/requests/respond", data);
    return response.data;
  },

  cancelFriendRequest: async (requestId: number) => {
    const response = await api.delete(`/friends/requests/${requestId}`);
    return response.data;
  },

  // Friendships
  getFriends: async () => {
    const response = await api.get("/friends/");
    return response.data;
  },

  removeFriend: async (friendId: number) => {
    const response = await api.delete(`/friends/${friendId}`);
    return response.data;
  },

  checkFriendship: async (friendId: number) => {
    const response = await api.get(`/friends/check/${friendId}`);
    return response.data;
  },

  // Blocked Users
  blockUser: async (data: { user_id: number; reason?: string }) => {
    const response = await api.post("/friends/block", data);
    return response.data;
  },

  unblockUser: async (userId: number) => {
    const response = await api.delete(`/friends/block/${userId}`);
    return response.data;
  },

  getBlockedUsers: async () => {
    const response = await api.get("/friends/blocked");
    return response.data;
  },

  checkUserBlocked: async (userId: number) => {
    const response = await api.get(`/friends/blocked/check/${userId}`);
    return response.data;
  },
};

export default api;
