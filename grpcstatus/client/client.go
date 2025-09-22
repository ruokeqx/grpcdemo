package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/ruokeqx/grpcdemo/grpcstatus/proto/tray/statuspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := statuspb.NewStatusServiceClient(conn)
	stream, err := client.StreamStatus(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	statuses := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}

	ticker := time.NewTicker(8 * time.Second)
	defer ticker.Stop()

	currentStatus := &statuspb.Status{
		UUID:      uuid.NewString(),
		TimeStamp: time.Now().Unix(),
		Msg:       "0",
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Receive error: %v", err)
				return
			}

			if msg.GetPullRequest() != nil {
				if err := stream.Send(&statuspb.StatusStreamMessage{
					Content: &statuspb.StatusStreamMessage_Status{Status: currentStatus},
				}); err != nil {
					log.Printf("Failed to send status on request: %v", err)
					return
				}
			}
		}
	}()

	for range ticker.C {
		newStatus := statuses[rand.Intn(len(statuses))]
		currentStatus = &statuspb.Status{
			UUID:      uuid.NewString(),
			TimeStamp: time.Now().Unix(),
			Msg:       newStatus,
		}

		log.Printf("Changing status to: %s", newStatus)

		if err := stream.Send(&statuspb.StatusStreamMessage{
			Content: &statuspb.StatusStreamMessage_Status{Status: currentStatus},
		}); err != nil {
			log.Printf("Failed to send status update: %v", err)
			return
		}
	}
}
