package services

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "grpc-location/proto/gen/go"
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

func (s *GreeterRPCService) LotsOfReplies(req *pb.HelloRequest, stream *pb.Greeter_LotsOfRepliesServer) error {
	return status.Errorf(codes.Unimplemented, "method LotsOfReplies not implemented")
}

func (s *GreeterRPCService) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
	return status.Errorf(codes.Unimplemented, "method LotsOfGreetings not implemented")
}

func (s *GreeterRPCService) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	return status.Errorf(codes.Unimplemented, "method BidiHello not implemented")
}
