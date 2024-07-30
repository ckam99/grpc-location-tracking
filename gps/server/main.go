package main

import (
	"google.golang.org/grpc/reflection"
	pb "grpc-location/proto/gen/go"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGPSServiceServer
	mu        sync.Mutex
	positions []*pb.Position
}

func (s *server) SendPosition(stream pb.GPSService_SendPositionServer) error {
	for {
		position, err := stream.Recv()
		if err != nil {
			return err
		}
		s.mu.Lock()
		s.positions = append(s.positions, position)
		s.mu.Unlock()
		stream.SendMsg(&pb.Response{Message: "Position received"})
	}
}

func (s *server) VisualizePosition(stream pb.GPSService_VisualizePositionServer) error {
	for {
		s.mu.Lock()
		for _, position := range s.positions {
			if err := stream.Send(position); err != nil {
				s.mu.Unlock()
				return err
			}
		}
		s.mu.Unlock()
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterGPSServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
