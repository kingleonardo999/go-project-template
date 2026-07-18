package main

import (
	"context"
	"log"
	"net"
	"net/http"

	helloPb "go-project-template/internal/api/v1/hello"
	"go-project-template/internal/config"
	helloService "go-project-template/internal/service/hello"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// StartGRPC 启动 gRPC 服务并注册所有服务
func StartGRPC(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("gRPC listen error: %v", err)
	}
	s := grpc.NewServer()
	helloPb.RegisterHelloServiceServer(s, helloService.NewService())
	reflection.Register(s)
	log.Printf("gRPC server listening on %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}

// StartHTTP 启动 HTTP 网关，将 HTTP 请求转发到 gRPC 后端
func StartHTTP(ctx context.Context, cfg *config.Config) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := helloPb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, cfg.Server.GRPC, opts); err != nil {
		log.Fatalf("register gateway error: %v", err)
	}

	srv := &http.Server{Addr: cfg.Server.HTTP, Handler: mux}
	log.Printf("HTTP gateway listening on %s", cfg.Server.HTTP)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP serve error: %v", err)
	}
}
