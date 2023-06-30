package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dspeirs7/animals/internal/api"
	"github.com/dspeirs7/animals/internal/log"
	"go.uber.org/zap"
)

func StartServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := log.NewLogger("server")
	defer logger.Sync()

	port := 8080

	api := api.NewAPI(ctx, logger)
	srv := api.Server(port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server couldn't be started", zap.Error(err))
		}
	}()

	logger.Info("server started", zap.Int("port", port))

	<-ctx.Done()
	stop()

	logger.Info("Starting gracefull shutdown")

	shutdownCtx, shutdownStop := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownStop()

	if err := api.Disconnect(shutdownCtx); err != nil {
		logger.Fatal("couldn't disconnect from db", zap.Error(err))
	}

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("gracefully shutdown")
}
