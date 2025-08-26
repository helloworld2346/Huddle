package config

import (
	"context"
	"fmt"
	"time"

	"huddle/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisClient *redis.Client

func InitRedis() error {
	config := GetConfig()
	
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: 10,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed to connect to Redis",
			zap.String("host", config.Redis.Host),
			zap.Int("port", config.Redis.Port),
			zap.Error(err),
		)
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	logger.Info("Redis connected successfully",
		zap.String("host", config.Redis.Host),
		zap.Int("port", config.Redis.Port),
		zap.Int("db", config.Redis.DB),
	)
	return nil
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func CloseRedis() error {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			logger.Error("Failed to close Redis connection", zap.Error(err))
			return err
		}
		logger.Info("Redis connection closed successfully")
	}
	return nil
}
