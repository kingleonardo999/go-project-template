package hello

import (
	"context"

	pb "go-project-template/internal/api/v1/hello"
)

type Service struct {
	ctx context.Context
	pb.UnimplementedHelloServiceServer
}

func NewService() *Service {
	return &Service{
		ctx: context.Background(),
	}
}

func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	return &pb.HelloResp{
		Message: "Hello, " + req.Name,
	}, nil
}
