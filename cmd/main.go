package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-hexagonal/adapter/dependency"
	"go-hexagonal/adapter/repository"
	"go-hexagonal/cmd/http_server"
	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

const ServiceName = "go-hexagonal"

// Constants for application settings
const (
	// DefaultShutdownTimeout is the default timeout for graceful shutdown
	DefaultShutdownTimeout = 5 * time.Second
)

func main() {
	fmt.Println("Starting " + ServiceName)

	// Initialize configuration
	config.Init("./config", "config")
	fmt.Println("Configuration initialized")

	// Initialize logging
	log.Init()
	fmt.Println("Logging initialized")

	// Create context and cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize repositories
	fmt.Println("Initializing repositories...")
	repository.Init(
		repository.WithMySQL(),
		// repository.WithRedis(), // Temporarily disabled Redis
	)
	fmt.Println("Repositories initialized")

	// Initialize services using dependency injection
	fmt.Println("Initializing services...")
	services, err := dependency.InitializeServices(ctx)
	if err != nil {
		log.SugaredLogger.Fatalf("Failed to initialize services: %v", err)
	}
	fmt.Println("Services initialized")

	// Create error channel and HTTP close channel
	errChan := make(chan error, 1)
	httpCloseCh := make(chan struct{}, 1)

	// Start HTTP server
	fmt.Println("Starting HTTP server...")
	go http_server.Start(ctx, errChan, httpCloseCh, services)
	fmt.Println("HTTP server started")

	// Listen for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal or error
	select {
	case err := <-errChan:
		log.SugaredLogger.Errorf("Server error: %v", err)
	case sig := <-sigChan:
		log.SugaredLogger.Infof("Received signal: %v", sig)
	}

	// Cancel context, trigger graceful shutdown
	log.SugaredLogger.Info("Shutting down server...")
	cancel()

	// Set shutdown timeout
	shutdownTimeout := DefaultShutdownTimeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	// Wait for HTTP server to close
	select {
	case <-httpCloseCh:
		log.SugaredLogger.Info("HTTP server shutdown completed")
	case <-shutdownCtx.Done():
		log.SugaredLogger.Warn("HTTP server shutdown timed out")
	}

	log.SugaredLogger.Info("Server gracefully stopped")
}
