package services

import (
	"google.golang.org/protobuf/types/known/emptypb"
	pb "grpc-location/proto/gen/go"
	"io"
	"sync"
)

type CarLocationService struct {
	pb.UnimplementedLocationServiceServer

	mu           sync.Mutex
	carLocations map[string][]*pb.CarLocation
}

func NewCarLocationService() *CarLocationService {
	return &CarLocationService{
		carLocations: make(map[string][]*pb.CarLocation),
	}
}

func (s *CarLocationService) SubmitCarLocation(stream pb.LocationService_SubmitCarLocationServer) error {
	for {
		carLocation, err := stream.Recv()
		if err == io.EOF {
			// End of the stream
			return nil
		}
		if err != nil {
			return err
		}

		s.mu.Lock()
		s.carLocations[carLocation.CarId] = append(s.carLocations[carLocation.CarId], carLocation)
		s.mu.Unlock()
	}
}

//func (s *CarLocationService) GetCarLocations(ctx context.Context, _ *emptypb.Empty) (pb.LocationService_GetCarLocationsServer, error) {
//
//	stream, err := grpc.NewServerStream(ctx, &pb.LocationService_GetCarLocationsServer{})
//	if err != nil {
//		return nil, err
//	}
//
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	for _, locations := range s.carLocations {
//		for _, location := range locations {
//			if err := stream.Send(location); err != nil {
//				return nil, err
//			}
//		}
//	}
//
//	return stream, nil
//}

func (s *CarLocationService) GetCarLocations(q *emptypb.Empty, stream pb.LocationService_GetCarLocationsServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, locations := range s.carLocations {
		for _, location := range locations {
			if err := stream.Send(location); err != nil {
				return err
			}
		}
	}

	return nil
}

// https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go
