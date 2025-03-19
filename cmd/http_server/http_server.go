package http_server

import (
	"context"
	"net/http"

	"github.com/spf13/cast"

	"go-hexagonal/adapter/dependency"
	http2 "go-hexagonal/api/http"
	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

// Start initializes and starts the HTTP server
func Start(ctx context.Context, errChan chan error, httpCloseCh chan struct{}) {
	// Initialize services using dependency injection
	_, err := dependency.InitializeServices(ctx)
	if err != nil {
		log.SugaredLogger.Errorf("Failed to initialize services: %v", err)
		errChan <- err
		return
	}

	// Initialize server
	srv := &http.Server{
		Addr:         config.GlobalConfig.HTTPServer.Addr,
		Handler:      http2.NewServerRoute(),
		ReadTimeout:  cast.ToDuration(config.GlobalConfig.HTTPServer.ReadTimeout),
		WriteTimeout: cast.ToDuration(config.GlobalConfig.HTTPServer.WriteTimeout),
	}

	// Run server
	go func() {
		log.SugaredLogger.Infof("%s HTTP server is starting on %s", config.GlobalConfig.App.Name, config.GlobalConfig.HTTPServer.Addr)
		errChan <- srv.ListenAndServe()
	}()

	// Watch for context cancellation
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.SugaredLogger.Infof("httpServer shutdown:%v", err)
		}
		httpCloseCh <- struct{}{}
	}()
}
