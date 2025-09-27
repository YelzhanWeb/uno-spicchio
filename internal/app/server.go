package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServerWithShutdown(addr string, mux *http.ServeMux) {
	logged := loggingMiddleware(mux)

	server := &http.Server{
		Addr:    addr,
		Handler: logged,
	}

	slog.Info("Starting server on", "port", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Start server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Shutdown error.")
	}
	slog.Info("Server gracefully stopped.")
}
