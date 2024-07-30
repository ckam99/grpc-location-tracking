package services

import (
	pb "grpc-location/proto/gen/go"
	"io"
	"log"
)

type ChatRPCService struct {
	pb.UnimplementedChatServiceServer
}

func NewChat() *ChatRPCService {
	return &ChatRPCService{}
}

func (s *ChatRPCService) ChatStream(stream pb.ChatService_ChatStreamServer) error {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return nil // End of stream
		}
		if err != nil {
			return err
		}
		log.Printf("Received message from %s: %s", message.GetUser(), message.GetMessage())

		// Echo the message back to the client
		err = stream.Send(&pb.ChatMessage{
			User:    "Server",
			Message: "Received: " + message.GetMessage(),
		})
		if err != nil {
			return err
		}
	}

}
