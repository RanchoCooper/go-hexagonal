package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go-hexagonal/cmd/http_server"
	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/internal/domain/service"
	"go-hexagonal/util/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	initConfig()
	initLogger()
	initRuntime(ctx)
	initServer(ctx, cancel)
}

func initConfig() {
	config.Init()
}

func initLogger() {
	log.Init()
}

func initRuntime(ctx context.Context) {
	repository.Init(
		repository.WithMySQL(),
		repository.WithRedis(),
	)
	service.Init(ctx)
}

func initServer(ctx context.Context, cancel context.CancelFunc) {
	errCh := make(chan error)
	httpCloseCh := make(chan struct{})
	http_server.Start(ctx, errCh, httpCloseCh)

	// graceful shutdown
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
	log.SugaredLogger.Infof("%s HTTP server exit!", config.Config.App.Name)
}
