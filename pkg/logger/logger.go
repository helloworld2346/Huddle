package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes the logger
func InitLogger() error {
	config := zap.NewProductionConfig()
	
	// Set log level
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	
	// Set encoding to JSON
	config.Encoding = "json"
	
	// Configure time format
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	
	// Configure level encoding
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	
	// Configure caller encoding
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	
	// Configure message key
	config.EncoderConfig.MessageKey = "message"
	
	// Set output paths
	config.OutputPaths = []string{"stdout", "logs/app.log"}
	config.ErrorOutputPaths = []string{"stderr", "logs/error.log"}
	
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}
	
	// Build logger
	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}
	
	return nil
}

// Info logs info level message
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Error logs error level message
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Warn logs warn level message
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Debug logs debug level message
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Fatal logs fatal level message and exits
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// WithContext creates a logger with request context
func WithContext(requestID, userID, method, path string) *zap.Logger {
	return Logger.With(
		zap.String("request_id", requestID),
		zap.String("user_id", userID),
		zap.String("method", method),
		zap.String("path", path),
	)
}

// HTTPRequest logs HTTP request details
func HTTPRequest(method, path, clientIP string, statusCode int, latency time.Duration, userAgent string) {
	Logger.Info("HTTP Request",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("client_ip", clientIP),
		zap.Int("status_code", statusCode),
		zap.Duration("latency", latency),
		zap.String("user_agent", userAgent),
	)
}

// DatabaseQuery logs database query details
func DatabaseQuery(query string, duration time.Duration, rows int) {
	Logger.Info("Database Query",
		zap.String("query", query),
		zap.Duration("duration", duration),
		zap.Int("rows", rows),
	)
}

// RedisOperation logs Redis operation details
func RedisOperation(operation, key string, duration time.Duration, success bool) {
	Logger.Info("Redis Operation",
		zap.String("operation", operation),
		zap.String("key", key),
		zap.Duration("duration", duration),
		zap.Bool("success", success),
	)
}

// Sync flushes any buffered log entries
func Sync() error {
	return Logger.Sync()
}
