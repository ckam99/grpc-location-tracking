package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc-location/proto/gen/go"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewLocationServiceClient(conn)

	// Submit car locations
	carLocations := []*pb.CarLocation{
		{CarId: "117k", Longitude: 10.0, Latitude: 20.0, Timestamp: uint64(time.Now().Unix())},
		{CarId: "1170", Longitude: 11.0, Latitude: 21.0, Timestamp: uint64(time.Now().Unix())},
	}

	stream, err := c.SubmitCarLocation(context.Background())
	if err != nil {
		log.Fatalf("could not submit car location: %v", err)
	}

	for _, location := range carLocations {
		if err := stream.Send(location); err != nil {
			log.Fatalf("could not send car location: %v", err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not close stream: %v", err)
	}

	// Get car locations
	resStream, err := c.GetCarLocation(context.Background(), &pb.CarLocationRequest{CarId: "636"})
	if err != nil {
		log.Fatalf("could not get car locations: %v", err)
	}

	for {
		carLocation, err := resStream.Recv()
		if err != nil {
			break
		}
		log.Printf("Received car location: %v", carLocation)
	}
}
