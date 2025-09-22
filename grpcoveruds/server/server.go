package main

import (
	"context"
	proto "demo/proto"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedEchoServer
}

func (s *server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	message := fmt.Sprintf("Hello, %s from server with PID %d", req.Name, os.Getpid())
	return &proto.HelloReply{Message: message}, nil
}

func main() {
	// 设置 UDS 地址
	socketPath := "/tmp/grpc_uds.sock"

	// 清理可能存在的旧socket文件
	if _, err := os.Stat(socketPath); err == nil {
		if err := os.Remove(socketPath); err != nil {
			log.Fatalf("Failed to remove existing socket: %v", err)
		}
	}

	// 创建 UDS 监听器
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	// 设置socket文件权限
	if err := os.Chmod(socketPath, 0666); err != nil {
		log.Fatalf("Failed to set socket permissions: %v", err)
	}

	// 创建 gRPC 服务器
	s := grpc.NewServer()
	proto.RegisterEchoServer(s, &server{})

	// 优雅关闭处理
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

	// 启动服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
