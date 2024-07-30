package main

import (
	"context"
	"fmt"
	pb "grpc-location/proto/gen/go"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func sendPosition(client pb.GPSServiceClient) {
	stream, err := client.SendPosition(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}
	for {
		position := &pb.Position{
			Latitude:  37.7749,
			Longitude: -122.4194,
			Timestamp: time.Now().Unix(),
		}
		if err := stream.Send(position); err != nil {
			log.Fatalf("Error sending position: %v", err)
		}
		// response, err := stream.Recv()
		response, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Error receiving response: %v", err)
		}
		log.Printf("Received response: %v", response.Message)
		time.Sleep(5 * time.Second)
	}
}

func visualizePosition(client pb.GPSServiceClient) {
	stream, err := client.VisualizePosition(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}
	for {
		position, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving position: %v", err)
		}
		log.Printf("Received position: %v, %v at %v", position.Latitude, position.Longitude, position.Timestamp)
	}
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGPSServiceClient(conn)

	// stop gracefully
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop // block console

	go sendPosition(client)
	//visualizePosition(client)

	fmt.Println("Shutting down...")

}
