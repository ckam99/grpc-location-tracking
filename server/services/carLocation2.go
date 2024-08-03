package services

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "grpc-location/proto/gen/go"
	"log"
	"sync"
)

type client struct {
	stream grpc.ServerStreamingServer[pb.CarLocation]
	car    string
}

func (c client) hasAttachedCar() bool {
	return c.car != ""
}

type empty struct{}

type CarLocationService2 struct {
	pb.UnimplementedLocationServiceServer
	mu              sync.Mutex
	carLocations    map[string]*pb.CarLocation
	locationChannel chan *pb.CarLocation
	quit            chan struct{}

	subscribers map[client]empty
	register    chan client
	unregister  chan client
}

func NewCarLocationService2() *CarLocationService2 {
	service := &CarLocationService2{
		carLocations: map[string]*pb.CarLocation{
			"967": {CarId: "967"},
			"895": {CarId: "895"},
			"288": {CarId: "288"},
			"636": {CarId: "636"},
		},
		locationChannel: make(chan *pb.CarLocation),
		quit:            make(chan struct{}),
		subscribers:     make(map[client]empty),
		register:        make(chan client),
		unregister:      make(chan client),
	}

	go service.processChannel()
	return service
}

func (s *CarLocationService2) processChannel() {
	for {
		select {
		case connection := <-s.register:
			s.subscribers[connection] = empty{}
			log.Println("new client registered")
		case car := <-s.locationChannel:
			s.mu.Lock()
			s.carLocations[car.CarId] = car
			s.mu.Unlock()
			for subscriber := range s.subscribers {
				if subscriber.hasAttachedCar() {
					if subscriber.car == car.CarId {
						if err := subscriber.stream.Send(car); err != nil {
							log.Printf("Failed to send car location: %v\n", err)
						}
					}
				} else if err := subscriber.stream.Send(car); err != nil {
					log.Printf("Failed to send car location: %v\n", err)
				}
			}
			log.Printf("Received and updated car location: %v", car)
		case connection := <-s.unregister:
			// Remove the client from the hub
			delete(s.subscribers, connection)
			log.Println("clients unregistered")

		case <-s.quit:
			close(s.locationChannel)
			return
		}
	}
}

func (s *CarLocationService2) SubmitLocation(ctx context.Context, car *pb.CarLocation) (*emptypb.Empty, error) {
	if _, err := s.findCarLocation(car.CarId); err != nil {
		return nil, err
	}
	go func() {
		s.mu.Lock()
		s.locationChannel <- car
		s.carLocations[car.CarId] = car
		s.mu.Unlock()
	}()
	return &emptypb.Empty{}, nil
}

func (s *CarLocationService2) GetAllCarLocation(_ *emptypb.Empty, stream grpc.ServerStreamingServer[pb.CarLocation]) error {
	for _, car := range s.carLocations {
		if car.Longitude > 0 {
			if err := stream.Send(car); err != nil {
				log.Printf("Failed to send car location: %v\n", err)
			}
		}
	}
	s.register <- client{stream: stream}
	// Wait for locationChannel closure signal
	<-s.quit
	return nil
}

func (s *CarLocationService2) GetCarLocation(req *pb.CarLocationRequest, stream grpc.ServerStreamingServer[pb.CarLocation]) error {
	log.Println("CarLocationRequest:", req)
	car, err := s.findCarLocation(req.CarId)
	if err != nil {
		return err
	}
	if car.Longitude > 0 {
		if err := stream.Send(car); err != nil {
			return status.Errorf(codes.Internal, "Failed to send car location: %v", err)
		}
	}
	s.register <- client{stream: stream, car: req.CarId}
	// Wait for locationChannel closure signal
	<-s.quit
	return nil
}

func (s *CarLocationService2) findCarLocation(carId string) (*pb.CarLocation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	car, ok := s.carLocations[carId]
	if !ok {
		return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("car location not found : %s", carId))
	}
	return car, nil
}

func (s *CarLocationService2) SubmitCarLocation(stream grpc.ClientStreamingServer[pb.CarLocation, emptypb.Empty]) error {
	return status.Errorf(codes.Unimplemented, "method SubmitLocation not implemented")
}

// https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go
