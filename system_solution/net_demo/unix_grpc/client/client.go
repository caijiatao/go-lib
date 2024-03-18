// client.go
package main

import (
	"context"
	"fmt"
	greater2 "golib/system_solution/net_demo/unix_grpc/pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// 连接到Unix域套接字
	conn, err := grpc.Dial("unix:///tmp/grpc_unix_socket", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建gRPC客户端
	client := greater2.NewGreeterClient(conn)

	// 发送gRPC请求
	response, err := client.SayHello(context.Background(), &greater2.HelloRequest{Name: "John"})
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	// 打印服务端的响应
	fmt.Println("Response from server:", response.Message)
}
