import { ApiResponse } from "./api";

export interface ValidationErrors {
  [key: string]: string;
}

export interface ErrorHandlerResult {
  fieldErrors: ValidationErrors;
  generalError?: string;
}

// Common backend error messages
const ERROR_MESSAGES = {
  USERNAME_EXISTS: "username already exists",
  EMAIL_EXISTS: "email already exists",
  INVALID_CREDENTIALS: "invalid credentials",
  USER_NOT_FOUND: "user not found",
  TOKEN_EXPIRED: "token expired",
  TOKEN_INVALID: "token invalid",
  PASSWORD_TOO_WEAK: "password too weak",
  EMAIL_INVALID: "invalid email format",
  USERNAME_INVALID: "invalid username format",
} as const;

// Map backend error messages to user-friendly messages
const ERROR_MESSAGE_MAP: Record<string, string> = {
  [ERROR_MESSAGES.USERNAME_EXISTS]: "Username already exists",
  [ERROR_MESSAGES.EMAIL_EXISTS]: "Email already exists",
  [ERROR_MESSAGES.INVALID_CREDENTIALS]: "Invalid username or password",
  [ERROR_MESSAGES.USER_NOT_FOUND]: "User not found",
  [ERROR_MESSAGES.TOKEN_EXPIRED]: "Session expired. Please login again",
  [ERROR_MESSAGES.TOKEN_INVALID]: "Invalid session. Please login again",
  [ERROR_MESSAGES.PASSWORD_TOO_WEAK]: "Password is too weak",
  [ERROR_MESSAGES.EMAIL_INVALID]: "Please enter a valid email address",
  [ERROR_MESSAGES.USERNAME_INVALID]: "Username contains invalid characters",
};

// Extract error message from API response
export function extractErrorMessage(response: ApiResponse): string {
  return response.error?.message || response.message || "An error occurred";
}

// Handle registration errors
export function handleRegistrationError(
  response: ApiResponse
): ErrorHandlerResult {
  const errorMessage = extractErrorMessage(response);
  const fieldErrors: ValidationErrors = {};
  let generalError: string | undefined;

  // Check for specific field errors
  if (errorMessage.includes(ERROR_MESSAGES.USERNAME_EXISTS)) {
    fieldErrors.username = ERROR_MESSAGE_MAP[ERROR_MESSAGES.USERNAME_EXISTS];
  } else if (errorMessage.includes(ERROR_MESSAGES.EMAIL_EXISTS)) {
    fieldErrors.email = ERROR_MESSAGE_MAP[ERROR_MESSAGES.EMAIL_EXISTS];
  } else if (errorMessage.includes(ERROR_MESSAGES.PASSWORD_TOO_WEAK)) {
    fieldErrors.password = ERROR_MESSAGE_MAP[ERROR_MESSAGES.PASSWORD_TOO_WEAK];
  } else if (errorMessage.includes(ERROR_MESSAGES.EMAIL_INVALID)) {
    fieldErrors.email = ERROR_MESSAGE_MAP[ERROR_MESSAGES.EMAIL_INVALID];
  } else if (errorMessage.includes(ERROR_MESSAGES.USERNAME_INVALID)) {
    fieldErrors.username = ERROR_MESSAGE_MAP[ERROR_MESSAGES.USERNAME_INVALID];
  } else {
    // General error
    generalError = errorMessage;
  }

  return { fieldErrors, generalError };
}

// Handle login errors
export function handleLoginError(response: ApiResponse): ErrorHandlerResult {
  const errorMessage = extractErrorMessage(response);
  const fieldErrors: ValidationErrors = {};
  let generalError: string | undefined;

  if (errorMessage.includes(ERROR_MESSAGES.INVALID_CREDENTIALS)) {
    generalError = ERROR_MESSAGE_MAP[ERROR_MESSAGES.INVALID_CREDENTIALS];
  } else if (errorMessage.includes(ERROR_MESSAGES.USER_NOT_FOUND)) {
    fieldErrors.username = ERROR_MESSAGE_MAP[ERROR_MESSAGES.USER_NOT_FOUND];
  } else {
    generalError = errorMessage;
  }

  return { fieldErrors, generalError };
}

// Handle general API errors
export function handleApiError(response: ApiResponse): ErrorHandlerResult {
  const errorMessage = extractErrorMessage(response);

  return {
    fieldErrors: {},
    generalError: errorMessage,
  };
}

// Network error handler
export function handleNetworkError(error: Error): ErrorHandlerResult {
  console.error("Network error:", error);

  return {
    fieldErrors: {},
    generalError: "Network error. Please check your connection and try again.",
  };
}

// Generic error handler
export function handleError(error: unknown): ErrorHandlerResult {
  if (error instanceof Error) {
    return handleNetworkError(error);
  }

  return {
    fieldErrors: {},
    generalError: "An unexpected error occurred. Please try again.",
  };
}
