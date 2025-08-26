package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"huddle/internal/config"
	"huddle/internal/database"
)

func main() {
	log.Println("🚀 Starting Huddle Server...")

	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("✅ Configuration loaded successfully")

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis
	if err := config.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	log.Println("🎉 All services initialized successfully!")
	log.Println("📊 Database: PostgreSQL connected")
	log.Println("🔴 Redis: Connected")
	log.Println("⚙️  Server ready to start...")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Close connections
	if err := database.CloseDatabase(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	if err := config.CloseRedis(); err != nil {
		log.Printf("Error closing Redis: %v", err)
	}

	log.Println("✅ Server shutdown complete")
}
