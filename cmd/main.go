package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-hexagonal/adapter/dependency"
	"go-hexagonal/adapter/repository"
	"go-hexagonal/cmd/http_server"
	"go-hexagonal/config"
	"go-hexagonal/domain/service"
	"go-hexagonal/util/log"
)

const ServiceName = "go-hexagonal"

func main() {
	fmt.Println("run " + ServiceName)

	ctx, cancel := context.WithCancel(context.Background())
	initConfig()
	initLogger()
	initRuntime(ctx)
	initServer(ctx, cancel)
}

func initConfig() {
	config.Init("./config", "config")
}

func initLogger() {
	log.Init()
}

func initRuntime(ctx context.Context) {
	// Initialize repositories
	repository.Init(
		repository.WithMySQL(),
		repository.WithRedis(),
	)

	// Initialize services using Wire dependency injection
	services, err := dependency.Injected(ctx)
	if err != nil {
		log.SugaredLogger.Fatalf("Failed to initialize services: %v", err)
	}

	// For backward compatibility, assign service instances to global variables
	service.ExampleSvc = services.ExampleService
	service.EventBus = services.EventBus
}

func initServer(ctx context.Context, cancel context.CancelFunc) {
	errCh := make(chan error)
	httpCloseCh := make(chan struct{})
	http_server.Start(ctx, errCh, httpCloseCh)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-quit:
		cancel()
		log.Logger.Info("Start graceful shutdown")
	case err := <-errCh:
		cancel()
		log.SugaredLogger.Errorf("http err:%v", err)
	}
	<-httpCloseCh
	log.SugaredLogger.Infof("%s HTTP server exit!", config.GlobalConfig.App.Name)
}
