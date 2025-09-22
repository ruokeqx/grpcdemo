package main

import (
	"context"
	proto "demo/proto"
	"fmt"

	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const DefaultIpcDialTimeout = 2 * time.Second

func main() {
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		log.Printf("Socket file does not exist: %s\n", socketPath)
	}

	conn, err := grpc.NewClient(
		targetPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(NewIpcConnection),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()

	fmt.Println("this is ok")

	client := proto.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := "World"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	response, err := client.SayHello(ctx, &proto.HelloRequest{Name: name})
	if err != nil {
		if se, ok := status.FromError(err); ok {
			if se.Code() != codes.OK {
				log.Fatalf("not ok: %v", se.Message())
			}
		}
		log.Fatalf("RPC failed: %v", err)
	}

	fmt.Printf("Server response: %s\n", response.Message)
}
