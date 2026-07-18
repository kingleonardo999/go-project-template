package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-project-template/internal/config"
)

func main() {
	cfg := config.NewConfig()
	if err := cfg.Load("config/app.yaml"); err != nil {
		log.Fatalf("config load error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartGRPC(cfg.Server.GRPC)
	go StartHTTP(ctx, cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
}
