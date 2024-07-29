package services

import (
	"context"
	pb "grpc-location/proto/gen"
	"log"
)

type GreeterRPCService struct {
	pb.UnimplementedGreeterServer
}

func NewGreeter() *GreeterRPCService {
	return &GreeterRPCService{}
}

// SayHello implements GreeterServer
func (s *GreeterRPCService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// SayHelloAgain implements GreeterServer
func (s *GreeterRPCService) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}
