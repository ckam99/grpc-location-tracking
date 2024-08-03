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

type CarLocationService struct {
	pb.UnimplementedLocationServiceServer
	mu                    sync.Mutex
	carLocations          map[string]*pb.CarLocation
	locationChannel       chan *pb.CarLocation
	singleLocationChannel chan *pb.CarLocation
	quit                  chan struct{}
}

func NewCarLocationService() *CarLocationService {
	return &CarLocationService{
		carLocations: map[string]*pb.CarLocation{
			"967": {CarId: "967", Longitude: 0.0, Latitude: 0.0},
			"895": {CarId: "895", Longitude: 0.0, Latitude: 0.0},
			"288": {CarId: "288", Longitude: 0.0, Latitude: 0.0},
			"636": {CarId: "636", Longitude: 0.0, Latitude: 0.0},
		},
		locationChannel:       make(chan *pb.CarLocation),
		singleLocationChannel: make(chan *pb.CarLocation),
		quit:                  make(chan struct{}),
	}
}

func (s *CarLocationService) SubmitCarLocation(stream grpc.ClientStreamingServer[pb.CarLocation, emptypb.Empty]) error {
	return status.Errorf(codes.Unimplemented, "method SubmitLocation not implemented")
}

func (s *CarLocationService) SubmitLocation(ctx context.Context, car *pb.CarLocation) (*emptypb.Empty, error) {
	if _, err := s.findCarLocation(car.CarId); err != nil {
		return nil, err
	}
	go func() {
		s.mu.Lock()
		s.locationChannel <- car
		s.singleLocationChannel <- car
		s.carLocations[car.CarId] = car
		s.mu.Unlock()
	}()
	return &emptypb.Empty{}, nil
}

func (s *CarLocationService) GetAllCarLocation(_ *emptypb.Empty, stream grpc.ServerStreamingServer[pb.CarLocation]) error {
	for _, car := range s.carLocations {
		if car.Longitude > 0 {
			if err := stream.Send(car); err != nil {
				log.Printf("Failed to send car location: %v\n", err)
			}
		}
	}
	for {
		select {
		case car := <-s.locationChannel:
			if err := stream.Send(car); err != nil {
				log.Printf("Failed to send car location: %v\n", err)
			}
		case <-s.quit:
			close(s.locationChannel)
			return nil
		}
		//car, ok := <-s.locationChannel
		//if !ok {
		//	return status.Errorf(codes.Unavailable, "service unavailable")
		//}
		//if err := stream.Send(car); err != nil {
		//	log.Printf("Failed to send car location: %v\n", err)
		//}
	}
}

func (s *CarLocationService) GetCarLocation(req *pb.CarLocationRequest, stream grpc.ServerStreamingServer[pb.CarLocation]) error {
	log.Println("CarLocationRequest:", req)
	if _, err := s.findCarLocation(req.CarId); err != nil {
		return err
	}
	for {
		select {
		case car := <-s.singleLocationChannel:
			if car.CarId == req.CarId {
				if err := stream.Send(car); err != nil {
					log.Printf("Failed to send car location: %v\n", err)
				}
			}
		case <-s.quit:
			close(s.singleLocationChannel)
			return nil
		}

		//car, ok := <-s.locationChannel
		//if !ok {
		//	return status.Errorf(codes.Unavailable, "service unavailable")
		//}
		//if car.CarId == req.CarId {
		//	if err := stream.Send(car); err != nil {
		//		log.Printf("Failed to send car location: %v\n", err)
		//	}
		//}
	}
}

func (s *CarLocationService) findCarLocation(carId string) (*pb.CarLocation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	car, ok := s.carLocations[carId]
	if !ok {
		return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("car location not found : %s", carId))
	}
	return car, nil
}

// https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go
