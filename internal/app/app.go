package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"huddle/internal/config"
	"huddle/internal/health"
	"huddle/internal/middleware"
	"huddle/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	router *gin.Engine
	server *http.Server
}

func NewApp() *App {
	// Set Gin mode
	if config.GetConfig().Server.Host == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.ErrorLogger())
	router.Use(gin.Recovery())

	// Setup routes
	setupRoutes(router)

	return &App{
		router: router,
	}
}

func setupRoutes(router *gin.Engine) {
	// API routes
	api := router.Group("/api")
	{
		// Health check routes (no auth required)
		health.SetupRoutes(api)
	}

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Huddle API",
			"version": "1.0.0",
			"docs":    "/api/health",
		})
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Route not found",
			"path":  c.Request.URL.Path,
		})
	})
}

func (a *App) Start() error {
	config := config.GetConfig()
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	a.server = &http.Server{
		Addr:         addr,
		Handler:      a.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("🚀 Server starting", zap.String("address", addr))
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.server != nil {
		logger.Info("🛑 Shutting down server...")
		return a.server.Shutdown(ctx)
	}
	return nil
}

func (a *App) GetRouter() *gin.Engine {
	return a.router
}
