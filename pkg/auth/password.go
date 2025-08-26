package auth

import (
	"errors"
	"unicode"

	"huddle/pkg/validation"

	"golang.org/x/crypto/bcrypt"
)

// Password requirements
const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
)

// Password validation rules - using validation package

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword checks if a password matches its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return errors.New("password must be at least 8 characters long")
	}
	
	if len(password) > MaxPasswordLength {
		return errors.New("password must be less than 128 characters")
	}
	
	if !validation.IsPasswordStrong(password) {
		return errors.New("password must contain uppercase, lowercase, number, and special character")
	}
	
	// Check for common weak passwords
	if validation.IsWeakPassword(password) {
		return errors.New("password is too common, please choose a stronger password")
	}
	
	return nil
}

// ValidateUsername validates username format
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	
	if len(username) > 20 {
		return errors.New("username must be less than 20 characters")
	}
	
	// Username can only contain alphanumeric characters and underscores
	if !validation.IsUsernameValid(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	
	// Username cannot start with a number
	if unicode.IsDigit(rune(username[0])) {
		return errors.New("username cannot start with a number")
	}
	
	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if !validation.IsEmailValid(email) {
		return errors.New("invalid email format")
	}
	
	if len(email) > 254 {
		return errors.New("email is too long")
	}
	
	return nil
}
