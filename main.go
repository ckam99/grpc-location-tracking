package main

import (
	"flag"
	"fmt"
	"grpc-location/server"
	"log"
	"log/slog"
	"os"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	serv := server.NewGrpcServer(slog.New(logHandler))
	if err := serv.Start(fmt.Sprintf(":%d", *port)); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
