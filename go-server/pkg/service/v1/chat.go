package v1

import (
	"../../api/v1"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"log"
)

// chatServiceServer is implementation of v1.ChatServiceServer proto interface
type chatServiceServer struct {
	msg chan string
}

// NewChatService Server creates Chat service object
func NewChatServiceServer() v1.ChatServiceServer {
	return &chatServiceServer{msg: make(chan string, 1000)}
}

// Send send message to the server
func (s *chatServiceServer) Send(ctx context.Context, message *wrappers.StringValue) (*empty.Empty, error) {
	if message != nil {
		log.Printf("Send requested: message=%v", *message)
		s.msg <- message.Value
	} else {
		log.Printf("Send requestedL message=<empty>")
	}

	return &empty.Empty{}, nil
}

// Subscribe is streaming method to get echo messages from the server
func (s *chatServiceServer) Subscribe(e *empty.Empty, stream v1.ChatService_SubscribeServer) error {
	log.Print("Subscribe requested")
	for {
		m := <- s.msg
		n := v1.Message{Text: fmt.Sprintf("I have received from you: %s. Thanks!", m)}
		if err := stream.Send(&n); err != nil {
			// put message back to the channel
			s.msg <- m
			log.Printf("Stream connection failed: %v", err)
			return nil
		}
		log.Printf("Message sent: %+v", n)
	}
}