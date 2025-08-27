"use client";

import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";
import { authAPI, userAPI } from "./api";

// User type
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

// Auth context type
interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<void>;
  register: (userData: RegisterData) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

// Register data type
export interface RegisterData {
  username: string;
  email: string;
  display_name: string;
  password: string;
}

// Create auth context
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Auth provider component
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Check if user is authenticated
  const isAuthenticated = !!user;

  // Initialize auth state
  useEffect(() => {
    const initAuth = async () => {
      try {
        const token = localStorage.getItem("access_token");
        if (token) {
          // Try to get current user
          const response = await userAPI.getCurrentUser();
          if (response.success && response.data) {
            setUser(response.data);
          } else {
            // Clear invalid tokens if user data is invalid
            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");
          }
        }
      } catch {
        // Clear invalid tokens
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
      } finally {
        setIsLoading(false);
      }
    };

    initAuth();
  }, []);

  // Login function
  const login = async (username: string, password: string) => {
    try {
      const response = await authAPI.login({ username, password });

      // Handle backend response format
      if (response.success && response.data) {
        const { tokens, user: userData } = response.data;

        // Store tokens
        localStorage.setItem("access_token", tokens.access_token);
        localStorage.setItem("refresh_token", tokens.refresh_token);

        // Set user
        setUser(userData);
      } else {
        throw new Error(response.error?.message || "Login failed");
      }
    } catch (error) {
      throw error;
    }
  };

  // Register function
  const register = async (userData: RegisterData) => {
    try {
      const response = await authAPI.register(userData);

      // Handle backend response format
      if (response.success && response.data) {
        const { tokens, user: newUser } = response.data;

        // Store tokens
        localStorage.setItem("access_token", tokens.access_token);
        localStorage.setItem("refresh_token", tokens.refresh_token);

        // Set user
        setUser(newUser);
      } else {
        throw new Error(response.error?.message || "Registration failed");
      }
    } catch (error) {
      throw error;
    }
  };

  // Logout function
  const logout = async () => {
    try {
      await authAPI.logout();
    } catch (error) {
      // Continue with logout even if API call fails
      console.error("Logout API error:", error);
    } finally {
      // Clear local state
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      setUser(null);
    }
  };

  // Refresh user data
  const refreshUser = async () => {
    try {
      const response = await userAPI.getCurrentUser();
      if (response.success && response.data) {
        setUser(response.data);
      } else {
        throw new Error("Failed to get user data");
      }
    } catch (error) {
      // If failed to get user, logout
      await logout();
      throw error;
    }
  };

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated,
    login,
    register,
    logout,
    refreshUser,
  };

  return React.createElement(AuthContext.Provider, { value }, children);
}

// Custom hook to use auth context
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
