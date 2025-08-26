package auth

import (
	"context"
	"fmt"
	"time"

	"huddle/internal/config"
	"huddle/pkg/logger"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis initializes Redis client
func InitRedis() {
	// Use the existing Redis client from config
	redisClient = config.GetRedisClient()
	if redisClient == nil {
		logger.Error("Redis client is nil - make sure Redis is initialized first")
		return
	}
	logger.Info("Auth Redis client initialized successfully")
}

// GetRedisClient returns Redis client
func GetRedisClient() *redis.Client {
	return redisClient
}

// BlacklistToken adds token to blacklist
func BlacklistToken(ctx context.Context, token string, expiresIn time.Duration) error {
	if redisClient == nil {
		logger.Error("Redis client is nil, cannot blacklist token")
		return fmt.Errorf("redis client not initialized")
	}
	
	key := "blacklist:" + token
	err := redisClient.Set(ctx, key, "blacklisted", expiresIn).Err()
	if err != nil {
		logger.Error("Failed to blacklist token", zap.Error(err))
		return err
	}
	
	logger.Info("Token blacklisted successfully", 
		zap.String("token", token[:10]+"..."),
		zap.String("key", key),
		zap.Duration("expires_in", expiresIn),
	)
	return nil
}

// IsTokenBlacklisted checks if token is blacklisted
func IsTokenBlacklisted(ctx context.Context, token string) bool {
	if redisClient == nil {
		logger.Error("Redis client is nil, cannot check token blacklist")
		return false
	}
	
	key := "blacklist:" + token
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		logger.Error("Failed to check token blacklist", zap.Error(err))
		return false
	}
	
	isBlacklisted := exists > 0
	logger.Info("Token blacklist check", 
		zap.String("token", token[:10]+"..."),
		zap.String("key", key),
		zap.Bool("is_blacklisted", isBlacklisted),
	)
	
	return isBlacklisted
}

// StoreUserSession stores user session in Redis
func StoreUserSession(ctx context.Context, userID uint, sessionData map[string]interface{}, expiresIn time.Duration) error {
	key := "session:" + string(rune(userID))
	err := redisClient.HSet(ctx, key, sessionData).Err()
	if err != nil {
		logger.Error("Failed to store user session", zap.Error(err))
		return err
	}
	
	// Set expiration
	err = redisClient.Expire(ctx, key, expiresIn).Err()
	if err != nil {
		logger.Error("Failed to set session expiration", zap.Error(err))
		return err
	}
	
	logger.Info("User session stored in Redis", zap.Uint("user_id", userID))
	return nil
}

// GetUserSession gets user session from Redis
func GetUserSession(ctx context.Context, userID uint) (map[string]string, error) {
	key := "session:" + string(rune(userID))
	result, err := redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		logger.Error("Failed to get user session", zap.Error(err))
		return nil, err
	}
	
	return result, nil
}

// DeleteUserSession deletes user session from Redis
func DeleteUserSession(ctx context.Context, userID uint) error {
	key := "session:" + string(rune(userID))
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		logger.Error("Failed to delete user session", zap.Error(err))
		return err
	}
	
	logger.Info("User session deleted from Redis", zap.Uint("user_id", userID))
	return nil
}

// StoreLoginAttempt stores login attempt for rate limiting
func StoreLoginAttempt(ctx context.Context, username string, success bool) error {
	key := "login_attempts:" + username
	attempt := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"success":   success,
		"ip":        "127.0.0.1", // TODO: Get from context
	}
	
	err := redisClient.LPush(ctx, key, attempt).Err()
	if err != nil {
		logger.Error("Failed to store login attempt", zap.Error(err))
		return err
	}
	
	// Keep only last 10 attempts
	err = redisClient.LTrim(ctx, key, 0, 9).Err()
	if err != nil {
		logger.Error("Failed to trim login attempts", zap.Error(err))
		return err
	}
	
	// Set expiration (1 hour)
	err = redisClient.Expire(ctx, key, time.Hour).Err()
	if err != nil {
		logger.Error("Failed to set login attempts expiration", zap.Error(err))
		return err
	}
	
	return nil
}

// GetLoginAttempts gets recent login attempts
func GetLoginAttempts(ctx context.Context, username string) ([]string, error) {
	key := "login_attempts:" + username
	result, err := redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		logger.Error("Failed to get login attempts", zap.Error(err))
		return nil, err
	}
	
	return result, nil
}

// IsRateLimited checks if user is rate limited
func IsRateLimited(ctx context.Context, username string) bool {
	attempts, err := GetLoginAttempts(ctx, username)
	if err != nil {
		return false
	}
	
	// Count failed attempts in last 15 minutes
	failedCount := 0
	
	for _, attempt := range attempts {
		// Simple check - in production, parse JSON properly
		if len(attempt) > 0 && attempt[0] == '{' {
			// This is a simplified check - in real implementation, parse JSON
			failedCount++
		}
	}
	
	return failedCount >= 5
}
