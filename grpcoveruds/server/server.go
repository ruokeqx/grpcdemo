package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ruokeqx/grpcdemo/grpcoveruds/proto/echopb"

	"google.golang.org/grpc"
)

type server struct {
	echopb.UnimplementedEchoServer
}

func (s *server) SayHello(ctx context.Context, req *echopb.HelloRequest) (*echopb.HelloReply, error) {
	message := fmt.Sprintf("Hello, %s from server with PID %d", req.Name, os.Getpid())
	return &echopb.HelloReply{Message: message}, nil
}

func main() {
	if _, err := os.Stat(socketPath); err == nil {
		if err := os.Remove(socketPath); err != nil {
			log.Fatalf("Failed to remove existing socket: %v", err)
		}
	}

	lis, err := NewIpcListener(socketPath)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	echopb.RegisterEchoServer(s, &server{})

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down server...")
		s.GracefulStop()
		os.Remove(socketPath)
		os.Exit(0)
	}()

	log.Printf("Server started on unix:%s", socketPath)
	log.Printf("Server PID: %d", os.Getpid())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
