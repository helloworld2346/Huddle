package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"huddle/internal/app"
	"huddle/internal/config"
	"huddle/internal/database"
	"huddle/pkg/auth"
	"huddle/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("🚀 Starting Huddle Server...")

	// Load configuration
	if err := config.Load(); err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}
	logger.Info("✅ Configuration loaded successfully")

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Initialize Redis
	if err := config.InitRedis(); err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}
	
	// Initialize Auth Redis
	auth.InitRedis()

	logger.Info("🎉 All services initialized successfully!")
	logger.Info("📊 Database: PostgreSQL connected")
	logger.Info("🔴 Redis: Connected")
	logger.Info("🔴 MinIO: Connected")

	// Create and start app
	app := app.NewApp()
	
	// Start app in a goroutine
	go func() {
		if err := app.Start(); err != nil {
			logger.Error("App error", zap.Error(err))
		}
	}()

	logger.Info("⚙️  Server ready to accept requests...")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown app
	if err := app.Shutdown(ctx); err != nil {
		logger.Error("App shutdown error", zap.Error(err))
	}

	// Close connections
	if err := database.CloseDatabase(); err != nil {
		logger.Error("Error closing database", zap.Error(err))
	}

	if err := config.CloseRedis(); err != nil {
		logger.Error("Error closing Redis", zap.Error(err))
	}

	logger.Info("✅ Server shutdown complete")
}
