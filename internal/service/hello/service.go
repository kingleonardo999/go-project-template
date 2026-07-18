package hello

import (
	"context"
	"log/slog"

	pb "go-project-template/internal/api/v1/hello"
)

type Service struct {
	pb.UnimplementedHelloServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	slog.InfoContext(ctx, "[hello] [SayHello]", "name", req.Name)
	return &pb.HelloResp{
		Message: "Hello, " + req.Name,
	}, nil
}
