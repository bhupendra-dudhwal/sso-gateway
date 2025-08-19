package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bhupendra-dudhwal/go-hexagonal/internal/builder"
	"go.uber.org/zap"
)

func main() {
	// Create base context
	ctx := context.Background()

	// Initialize application via builder pattern
	logger, server, port := builder.
		NewAppBuilder(ctx).
		SetConfig().
		SetLogger().
		SetDatabase().
		SetCache().
		SetServices().
		SetHandler().
		Build()

	logger.Info("Application builder completed successfully")
	logger.Info("Starting server", zap.Int("port", port))

	// Start the HTTP server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%d", port)
		if err := server.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
			logger.Error("Server startup error", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until signal is received
	sig := <-stop
	logger.Info("Shutdown signal received", zap.Any("signal", sig))

	// Context for graceful shutdown
	shutdownCtx, cancel := context.WithTimeoutCause(
		ctx,
		10*time.Second,
		errors.New("server interrupt by os signal"),
	)
	defer cancel()

	if err := server.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Info("server shutdown error", zap.Error(err))
	} else {
		logger.Info("server gracefully stopped")
	}
}
