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

	"go.uber.org/zap"
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
	log.Logger.Info("Application starting",
		zap.String("service", ServiceName),
		zap.String("env", string(config.GlobalConfig.Env)))

	// Create context and cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize repositories using wire dependency injection with options
	log.Logger.Info("Initializing repositories")
	clients, err := dependency.InitializeRepositories(
		dependency.WithMySQL(),
		dependency.WithRedis(),
	)
	if err != nil {
		log.Logger.Fatal("Failed to initialize repositories",
			zap.Error(err))
	}
	repository.Clients = clients
	log.Logger.Info("Repositories initialized successfully")

	// Initialize services using dependency injection
	log.Logger.Info("Initializing services")
	services, err := dependency.InitializeServices(ctx, dependency.WithExampleService())
	if err != nil {
		log.Logger.Fatal("Failed to initialize services",
			zap.Error(err))
	}
	log.Logger.Info("Services initialized successfully")

	// Create error channel and HTTP close channel
	errChan := make(chan error, 1)
	httpCloseCh := make(chan struct{}, 1)

	// Start HTTP server
	log.Logger.Info("Starting HTTP server",
		zap.String("address", config.GlobalConfig.HTTPServer.Addr))
	go http_server.Start(ctx, errChan, httpCloseCh, services)
	log.Logger.Info("HTTP server started")

	// Listen for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal or error
	select {
	case err := <-errChan:
		log.Logger.Error("Server error", zap.Error(err))
	case sig := <-sigChan:
		log.Logger.Info("Received signal", zap.String("signal", sig.String()))
	}

	// Cancel context, trigger graceful shutdown
	log.Logger.Info("Shutting down server")
	cancel()

	// Set shutdown timeout
	shutdownTimeout := DefaultShutdownTimeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	// Wait for HTTP server to close
	select {
	case <-httpCloseCh:
		log.Logger.Info("HTTP server shutdown completed")
	case <-shutdownCtx.Done():
		log.Logger.Warn("HTTP server shutdown timed out",
			zap.Duration("timeout", DefaultShutdownTimeout))
	}

	log.Logger.Info("Server gracefully stopped")
}
