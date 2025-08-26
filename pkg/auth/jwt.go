package auth

import (
	"errors"
	"time"

	"huddle/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// GenerateTokenPair generates access and refresh tokens
func GenerateTokenPair(userID uint, username, email string) (*TokenPair, error) {
	config := config.GetConfig()
	
	// Access token - 15 minutes
	accessExp := time.Now().Add(15 * time.Minute)
	accessClaims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "huddle",
			Subject:   username,
		},
	}
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(config.JWT.Secret))
	if err != nil {
		return nil, err
	}
	
	// Refresh token - 7 days
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "huddle",
			Subject:   username,
		},
	}
	
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.JWT.Secret))
	if err != nil {
		return nil, err
	}
	
	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    900, // 15 minutes in seconds
	}, nil
}

// ValidateToken validates and parses JWT token
func ValidateToken(tokenString string) (*JWTClaims, error) {
	config := config.GetConfig()
	
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.JWT.Secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

// RefreshToken generates new token pair from refresh token
func RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := ValidateToken(refreshTokenString)
	if err != nil {
		return nil, err
	}
	
	// Generate new token pair
	return GenerateTokenPair(claims.UserID, claims.Username, claims.Email)
}

// ExtractUserIDFromToken extracts user ID from token without full validation
func ExtractUserIDFromToken(tokenString string) (uint, error) {
	config := config.GetConfig()
	
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Secret), nil
	})
	
	if err != nil {
		return 0, err
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims.UserID, nil
	}
	
	return 0, errors.New("invalid token")
}
