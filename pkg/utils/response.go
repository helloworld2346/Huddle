package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response
type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorInfo represents error details
type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	response := Response{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}
	c.JSON(http.StatusOK, response)
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, code, message string, details map[string]interface{}) {
	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}
	c.JSON(statusCode, response)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, details map[string]interface{}) {
	ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
}

// UnauthorizedResponse sends an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// ForbiddenResponse sends a forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", message, nil)
}

// NotFoundResponse sends a not found response
func NotFoundResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", message, nil)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

// BadRequestResponse sends a bad request response
func BadRequestResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Bad request"
	}
	ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

// ConflictResponse sends a conflict response
func ConflictResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource conflict"
	}
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message, nil)
}

// TooManyRequestsResponse sends a rate limit response
func TooManyRequestsResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Too many requests"
	}
	ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT", message, nil)
}
