// server.go
package main

import (
	"context"
	"fmt"
	"golib/examples/net_demo/unix_grpc/pb"
	"net"
	"os"

	"google.golang.org/grpc"
)

type server struct {
	greater.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *greater.HelloRequest) (*greater.HelloReply, error) {
	return &greater.HelloReply{Message: "Hello, " + req.Name}, nil
}

func main() {
	// 创建Unix域套接字
	sockPath := "/tmp/grpc_unix_socket"
	os.Remove(sockPath)

	lis, err := net.Listen("unix", sockPath)
	if err != nil {
		fmt.Printf("Failed to create listener: %v", err)
		return
	}
	defer os.Remove(sockPath)

	// 创建gRPC服务器
	srv := grpc.NewServer()
	greater.RegisterGreeterServer(srv, &server{})

	// 启动gRPC服务器
	fmt.Println("gRPC server is running on Unix socket:", sockPath)
	if err := srv.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
