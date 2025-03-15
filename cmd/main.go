package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-hexagonal/adapter/repository"
	"go-hexagonal/cmd/http_server"
	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

const ServiceName = "go-hexagonal"

func main() {
	fmt.Println("run " + ServiceName)

	// Initialize configuration
	config.Init("./config", "config")

	// Initialize logging
	log.Init()

	// Create context and cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize repositories
	repository.Init(
		repository.WithMySQL(),
		repository.WithRedis(),
	)
<<<<<<< HEAD
}
=======

	// Create error channel and HTTP close channel
	errChan := make(chan error, 1)
	httpCloseCh := make(chan struct{}, 1)

	// Start HTTP server
	go http_server.Start(ctx, errChan, httpCloseCh)
>>>>>>> 4821dda (chore: code improve)

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
	shutdownTimeout := 5 * time.Second
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	// Wait for HTTP server to close or timeout
	select {
	case <-httpCloseCh:
		log.SugaredLogger.Info("HTTP server shutdown complete")
	case <-shutdownCtx.Done():
		log.SugaredLogger.Warn("HTTP server shutdown timed out")
	}

	log.SugaredLogger.Info("Server exited")
}
