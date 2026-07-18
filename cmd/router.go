package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"

	"go-project-template/internal/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// StartGRPC 启动 gRPC 服务并注册所有服务
func StartGRPC(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("gRPC listen error", "error", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	registerGRPCServices(s)
	reflection.Register(s)
	slog.Info("gRPC server listening", "addr", addr)
	if err := s.Serve(lis); err != nil {
		slog.Error("gRPC serve error", "error", err)
		os.Exit(1)
	}
}

// StartHTTP 启动 HTTP 网关，将 HTTP 请求转发到 gRPC 后端
func StartHTTP(ctx context.Context, cfg *config.Config) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := registerHTTPHandlers(ctx, mux, cfg.Server.GRPC, opts); err != nil {
		slog.Error("register gateway error", "error", err)
		os.Exit(1)
	}

	srv := &http.Server{Addr: cfg.Server.HTTP, Handler: mux}
	slog.Info("HTTP gateway listening", "addr", cfg.Server.HTTP)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP serve error", "error", err)
		os.Exit(1)
	}
}
