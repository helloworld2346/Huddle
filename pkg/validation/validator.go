package validation

import (
	"strings"
	"unicode"

	"huddle/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// AddError adds a validation error
func (ve *ValidationErrors) AddError(field, message string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors checks if there are any validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// ToMap converts validation errors to a map for response
func (ve *ValidationErrors) ToMap() map[string]interface{} {
	errors := make(map[string]interface{})
	for _, err := range ve.Errors {
		errors[err.Field] = err.Message
	}
	return errors
}

// ValidateRequired validates required fields
func ValidateRequired(value, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " is required",
		}
	}
	return nil
}

// ValidateLength validates string length
func ValidateLength(value, fieldName string, min, max int) *ValidationError {
	length := len(strings.TrimSpace(value))
	if length < min {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be at least " + string(rune(min)) + " characters long",
		}
	}
	if max > 0 && length > max {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be less than " + string(rune(max)) + " characters",
		}
	}
	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) *ValidationError {
	if email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "email is required",
		}
	}
	
	if !IsEmailValid(email) {
		return &ValidationError{
			Field:   "email",
			Message: "invalid email format",
		}
	}
	
	if len(email) > 254 {
		return &ValidationError{
			Field:   "email",
			Message: "email is too long",
		}
	}
	
	return nil
}

// ValidateUsername validates username format
func ValidateUsername(username string) *ValidationError {
	if username == "" {
		return &ValidationError{
			Field:   "username",
			Message: "username is required",
		}
	}
	
	if len(username) < 3 {
		return &ValidationError{
			Field:   "username",
			Message: "username must be at least 3 characters long",
		}
	}
	
	if len(username) > 20 {
		return &ValidationError{
			Field:   "username",
			Message: "username must be less than 20 characters",
		}
	}
	
	// Username can only contain alphanumeric characters and underscores
	if !IsUsernameValid(username) {
		return &ValidationError{
			Field:   "username",
			Message: "username can only contain letters, numbers, and underscores",
		}
	}
	
	// Username cannot start with a number
	if unicode.IsDigit(rune(username[0])) {
		return &ValidationError{
			Field:   "username",
			Message: "username cannot start with a number",
		}
	}
	
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) *ValidationError {
	if password == "" {
		return &ValidationError{
			Field:   "password",
			Message: "password is required",
		}
	}
	
	if len(password) < 8 {
		return &ValidationError{
			Field:   "password",
			Message: "password must be at least 8 characters long",
		}
	}
	
	if len(password) > 128 {
		return &ValidationError{
			Field:   "password",
			Message: "password must be less than 128 characters",
		}
	}
	
	if !IsPasswordStrong(password) {
		return &ValidationError{
			Field:   "password",
			Message: "password must contain uppercase, lowercase, number, and special character",
		}
	}
	
	// Check for common weak passwords
	if IsWeakPassword(password) {
		return &ValidationError{
			Field:   "password",
			Message: "password is too common, please choose a stronger password",
		}
	}
	
	return nil
}

// ValidateConfirmPassword validates password confirmation
func ValidateConfirmPassword(password, confirmPassword string) *ValidationError {
	if confirmPassword == "" {
		return &ValidationError{
			Field:   "confirm_password",
			Message: "password confirmation is required",
		}
	}
	
	if password != confirmPassword {
		return &ValidationError{
			Field:   "confirm_password",
			Message: "passwords do not match",
		}
	}
	
	return nil
}

// ValidateBio validates bio length
func ValidateBio(bio string) *ValidationError {
	if bio != "" && len(bio) > 500 {
		return &ValidationError{
			Field:   "bio",
			Message: "bio must be less than 500 characters",
		}
	}
	return nil
}

// ValidateDisplayName validates display name
func ValidateDisplayName(displayName string) *ValidationError {
	if displayName == "" {
		return &ValidationError{
			Field:   "display_name",
			Message: "display name is required",
		}
	}
	
	if len(displayName) < 2 {
		return &ValidationError{
			Field:   "display_name",
			Message: "display name must be at least 2 characters long",
		}
	}
	
	if len(displayName) > 50 {
		return &ValidationError{
			Field:   "display_name",
			Message: "display name must be less than 50 characters",
		}
	}
	
	return nil
}

// SendValidationError sends validation error response
func SendValidationError(c *gin.Context, errors *ValidationErrors) {
	utils.ValidationErrorResponse(c, errors.ToMap())
}

// ValidateAndSendError validates and sends error if validation fails
func ValidateAndSendError(c *gin.Context, errors *ValidationErrors) bool {
	if errors.HasErrors() {
		SendValidationError(c, errors)
		return true
	}
	return false
}
