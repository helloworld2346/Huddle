package database

import (
	"fmt"
	"time"

	"huddle/internal/config"
	"huddle/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() error {
	config := config.GetConfig()
	
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})

	if err != nil {
		logger.Error("Failed to connect to database",
			zap.String("host", config.Database.Host),
			zap.Int("port", config.Database.Port),
			zap.String("database", config.Database.Name),
			zap.Error(err),
		)
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Error("Failed to get underlying sql.DB", zap.Error(err))
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return fmt.Errorf("failed to ping database: %v", err)
	}

	logger.Info("Database connected successfully",
		zap.String("host", config.Database.Host),
		zap.Int("port", config.Database.Port),
		zap.String("database", config.Database.Name),
	)
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			logger.Error("Failed to get sql.DB for closing", zap.Error(err))
			return err
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
			return err
		}
		logger.Info("Database connection closed successfully")
	}
	return nil
}
