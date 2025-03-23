// Package log provides logging functionality for the application
package log

import (
	"errors"
	"time"

	"go.uber.org/zap"
)

// LoggerUsageExamples demonstrates the best practices for using Logger and SugaredLogger
func LoggerUsageExamples() {
	// Sample variables for demonstration
	errDatabaseConnection := errors.New("connection refused")
	userID := "user456"
	ipAddress := "203.0.113.42"
	err := errors.New("invalid request parameters")

	// ===== LOGGER (Structured Logger) =====
	// Best used for:
	// 1. Performance-critical code paths
	// 2. Structured data with known fields
	// 3. High-volume logging

	// Example 1: Basic structured logging with explicit fields
	Logger.Info("User logged in",
		zap.String("user_id", "user123"),
		zap.String("ip_address", "192.168.1.1"),
		zap.String("user_agent", "Mozilla/5.0"),
	)

	// Example 2: Error logging with structured context
	Logger.Error("Database connection failed",
		zap.String("db_host", "db.example.com"),
		zap.Int("port", 5432),
		zap.Duration("timeout", 30*time.Second),
		zap.Error(errDatabaseConnection),
	)

	// Example 3: Warn level with context fields
	Logger.Warn("Rate limit exceeded",
		zap.String("client_id", "client456"),
		zap.Int("limit", 100),
		zap.Int("current_rate", 120),
	)

	// ===== SUGARED LOGGER =====
	// Best used for:
	// 1. Debug/development logging
	// 2. Dynamic or variable number of arguments
	// 3. When convenience is more important than optimal performance

	// Example 1: Simple message with printf-style formatting
	SugaredLogger.Infof("Processing item %d of %d", 5, 10)

	// Example 2: Error with variable message
	SugaredLogger.Errorf("Failed to process request: %v", err)

	// Example 3: Using key-value pairs with .Infow, .Errorw, etc.
	SugaredLogger.Infow("Request processed",
		"method", "GET",
		"path", "/api/users",
		"status", 200,
		"duration_ms", 45.2,
	)

	// Example 4: Warning with formatting
	SugaredLogger.Warnf("Unusual access pattern detected for user %s from IP %s", userID, ipAddress)
}

// General recommendation on when to use each:
//
// Use Logger (structured logger) when:
// - You're logging in a hot code path (performance critical)
// - You have a fixed set of fields to log
// - You need maximum performance
// - You're building core infrastructure or libraries
//
// Use SugaredLogger when:
// - You're logging in non-performance-critical paths
// - You need string formatting
// - You have a variable number of fields
// - You're writing application code where convenience is important
// - You're doing temporary debugging
