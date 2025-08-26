package validation

import (
	"regexp"
)

// Regex patterns for validation
var (
	// Email validation: user@domain.com
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	// Username validation: alphanumeric + underscore, 3-20 chars
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	
	// Password strength validation
	HasUpperRegex   = regexp.MustCompile(`[A-Z]`)
	HasLowerRegex   = regexp.MustCompile(`[a-z]`)
	HasNumberRegex  = regexp.MustCompile(`[0-9]`)
	HasSpecialRegex = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	
	// Display name validation: letters, spaces, dots, hyphens
	DisplayNameRegex = regexp.MustCompile(`^[a-zA-Z\s\.\-]+$`)
	
	// Phone number validation (basic)
	PhoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	
	// URL validation
	URLRegex = regexp.MustCompile(`^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`)
	
	// UUID validation
	UUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	
	// Date format validation (YYYY-MM-DD)
	DateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	
	// Time format validation (HH:MM:SS)
	TimeRegex = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`)
	
	// DateTime format validation (YYYY-MM-DD HH:MM:SS)
	DateTimeRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	
	// IP address validation
	IPRegex = regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	
	// Hex color validation (#RRGGBB or #RRGGBBAA)
	HexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{8})$`)
	
	// Base64 validation
	Base64Regex = regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
	
	// Alphanumeric validation
	AlphanumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	
	// Letters only validation
	LettersOnlyRegex = regexp.MustCompile(`^[a-zA-Z]+$`)
	
	// Numbers only validation
	NumbersOnlyRegex = regexp.MustCompile(`^[0-9]+$`)
	
	// No special characters validation (except spaces)
	NoSpecialCharsRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
)

// Common weak passwords to check against
var WeakPasswords = []string{
	"password", "123456", "qwerty", "admin", "letmein",
	"welcome", "monkey", "dragon", "master", "hello",
	"123456789", "12345678", "1234567", "1234567890",
	"abc123", "password123", "admin123", "root",
	"guest", "user", "test", "demo", "sample",
}

// IsEmailValid checks if email format is valid
func IsEmailValid(email string) bool {
	return EmailRegex.MatchString(email)
}

// IsUsernameValid checks if username format is valid
func IsUsernameValid(username string) bool {
	return UsernameRegex.MatchString(username)
}

// IsPasswordStrong checks if password meets strength requirements
func IsPasswordStrong(password string) bool {
	return HasUpperRegex.MatchString(password) &&
		HasLowerRegex.MatchString(password) &&
		HasNumberRegex.MatchString(password) &&
		HasSpecialRegex.MatchString(password)
}

// IsWeakPassword checks if password is in common weak passwords list
func IsWeakPassword(password string) bool {
	for _, weak := range WeakPasswords {
		if password == weak {
			return true
		}
	}
	return false
}

// IsDisplayNameValid checks if display name format is valid
func IsDisplayNameValid(displayName string) bool {
	return DisplayNameRegex.MatchString(displayName)
}

// IsPhoneValid checks if phone number format is valid
func IsPhoneValid(phone string) bool {
	return PhoneRegex.MatchString(phone)
}

// IsURLValid checks if URL format is valid
func IsURLValid(url string) bool {
	return URLRegex.MatchString(url)
}

// IsUUIDValid checks if UUID format is valid
func IsUUIDValid(uuid string) bool {
	return UUIDRegex.MatchString(uuid)
}

// IsDateValid checks if date format is valid (YYYY-MM-DD)
func IsDateValid(date string) bool {
	return DateRegex.MatchString(date)
}

// IsTimeValid checks if time format is valid (HH:MM:SS)
func IsTimeValid(time string) bool {
	return TimeRegex.MatchString(time)
}

// IsDateTimeValid checks if datetime format is valid
func IsDateTimeValid(datetime string) bool {
	return DateTimeRegex.MatchString(datetime)
}

// IsIPValid checks if IP address format is valid
func IsIPValid(ip string) bool {
	return IPRegex.MatchString(ip)
}

// IsHexColorValid checks if hex color format is valid
func IsHexColorValid(color string) bool {
	return HexColorRegex.MatchString(color)
}

// IsBase64Valid checks if string is valid base64
func IsBase64Valid(str string) bool {
	return Base64Regex.MatchString(str)
}

// IsAlphanumeric checks if string contains only alphanumeric characters
func IsAlphanumeric(str string) bool {
	return AlphanumericRegex.MatchString(str)
}

// IsLettersOnly checks if string contains only letters
func IsLettersOnly(str string) bool {
	return LettersOnlyRegex.MatchString(str)
}

// IsNumbersOnly checks if string contains only numbers
func IsNumbersOnly(str string) bool {
	return NumbersOnlyRegex.MatchString(str)
}

// IsNoSpecialChars checks if string contains no special characters (except spaces)
func IsNoSpecialChars(str string) bool {
	return NoSpecialCharsRegex.MatchString(str)
}
