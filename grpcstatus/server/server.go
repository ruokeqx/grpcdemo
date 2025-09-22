package main

import (
	"log"
	"net"
	"time"

	"github.com/ruokeqx/grpcdemo/grpcstatus/proto/tray/statuspb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	statuspb.UnimplementedStatusServiceServer
	currentStatus *statuspb.Status
}

func (s *server) StreamStatus(stream statuspb.StatusService_StreamStatusServer) error {
	log.Println("Client connected")
	defer log.Println("Client disconnected")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	msgChan := make(chan *statuspb.StatusStreamMessage, 10)

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Receive error: %v", err)
				close(msgChan)
				return
			}
			msgChan <- msg
		}
	}()

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			log.Printf("get msg: %v", msg)

		case <-ticker.C:
			if err := stream.Send(&statuspb.StatusStreamMessage{
				Content: &statuspb.StatusStreamMessage_PullRequest{
					PullRequest: &emptypb.Empty{},
				},
			}); err != nil {
				log.Printf("Send pull request error: %v", err)
				return err
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	statuspb.RegisterStatusServiceServer(s, &server{})

	log.Println("Server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
