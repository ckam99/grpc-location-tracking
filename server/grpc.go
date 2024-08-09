package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"grpc-location/docs"
	v1 "grpc-location/proto/gen/go"
	"grpc-location/server/interceptor"
	"grpc-location/server/services"
	"log/slog"
	"net"
	"net/http"
	"path/filepath"
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
	//v1.RegisterLocationServiceServer(s.grpcServer, services.NewCarLocationService())
	v1.RegisterLocationServiceServer(s.grpcServer, services.NewCarLocationService2())

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("Failed to listen: %w", err)
	}
	defer listener.Close()
	// s.log.Info(fmt.Sprintf("GRPC server is running on port %s.", address), "op", "Start")
	go func() {
		if err := runGateway(address, s.log); err != nil {
			s.log.Error("Failed to start gateway", "error", err)
		}
	}()
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

func runGateway(addr string, log *slog.Logger) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	rmux := runtime.NewServeMux()
	opts := []gRPC.DialOption{gRPC.WithTransportCredentials(insecure.NewCredentials())}
	err := v1.RegisterLocationServiceHandlerFromEndpoint(ctx, rmux, addr, opts)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(filepath.Join("./", "docs/swagger")))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	mux.Handle("/", rmux)
	mux.Handle("/docs/swagger", docs.SwaggerHandler)

	log.Info(fmt.Sprintf("server listening at %v", "http://0.0.0.0:8081"))
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}
