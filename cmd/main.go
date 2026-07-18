package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-project-template/internal/config"
)

func main() {
	cfg := config.NewConfig()
	if err := cfg.Load("config/app.yaml"); err != nil {
		slog.Error("config load error", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartGRPC(cfg.Server.GRPC)
	go StartHTTP(ctx, cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down...")
}
