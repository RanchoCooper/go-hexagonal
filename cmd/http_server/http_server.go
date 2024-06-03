package http_server

import (
	"context"
	"net/http"

	"github.com/spf13/cast"

	http2 "go-hexagonal/api/http"
	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

func Start(ctx context.Context, errChan chan error, httpCloseCh chan struct{}) {
	// init server
	srv := &http.Server{
		Addr:         config.Config.HTTPServer.Addr,
		Handler:      http2.NewServerRoute(),
		ReadTimeout:  cast.ToDuration(config.Config.HTTPServer.ReadTimeout),
		WriteTimeout: cast.ToDuration(config.Config.HTTPServer.WriteTimeout),
	}

	// run server
	go func() {
		log.SugaredLogger.Infof("%s HTTP server is starting on %s", config.Config.App.Name, config.Config.HTTPServer.Addr)
		errChan <- srv.ListenAndServe()
	}()

	// watch the ctx exit
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.SugaredLogger.Infof("httpServer shutdown:%v", err)
		}
		httpCloseCh <- struct{}{}
	}()
}
