package server

import (
	"fmt"
	gRPC "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	v1 "grpc-location/proto/gen/go"
	"grpc-location/server/interceptor"
	"grpc-location/server/services"
	"log/slog"
	"net"
)

type Server struct {
	grpcServer *gRPC.Server
	log        *slog.Logger
}

func NewGrpcServer(logger *slog.Logger) *Server {
	return &Server{
		log: logger,
	}
}

func (s *Server) Start(address string) error {
	rpcLogger := gRPC.UnaryInterceptor(interceptor.GrpcUnaryLogger(s.log))
	s.grpcServer = gRPC.NewServer(rpcLogger)
	reflection.Register(s.grpcServer)
	//v1.RegisterGreeterServer(s.grpcServer, services.NewGreeter())
	v1.RegisterChatServiceServer(s.grpcServer, services.NewChat())
	v1.RegisterLocationServiceServer(s.grpcServer, services.NewCarLocationService())

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("Failed to listen: %w", err)
	}
	defer listener.Close()
	// s.log.Info(fmt.Sprintf("GRPC server is running on port %s.", address), "op", "Start")
	s.log.Info(fmt.Sprintf("server listening at %v", listener.Addr()))
	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("cannot start gRPC server: %w", err)
	}
	return nil
}

func (s *Server) Close() error {
	s.grpcServer.GracefulStop()
	return nil
}
