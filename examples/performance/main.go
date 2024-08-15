package main

import (
	"context"
	"golib/examples/performance/feature"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const (
	address       = "localhost:8900" // gRPC 服务地址
	concurrency   = 100              // 并发请求数
	totalRequests = 10000            // 总请求数
)

func main() {
	// 创建 gRPC 连接
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := feature.NewFeatureServerClient(conn) // 替换为你的服务客户端

	var g errgroup.Group
	var mu sync.Mutex
	successCount := 0

	startTime := time.Now()

	// 并发执行 gRPC 请求
	for i := 0; i < concurrency; i++ {
		g.Go(func() error {
			for j := 0; j < totalRequests/concurrency; j++ {
				_, err := client.GetUserPublishedPaper(context.Background(), &feature.GetUserPublishedPaperRequest{
					UserId: "test-scholar-topic-1-001",
				})
				if err != nil {
					log.Printf("Request failed: %v", err)
					return err
				}
				mu.Lock()
				successCount++
				mu.Unlock()
			}
			return nil
		})
	}

	// 等待所有请求完成
	if err := g.Wait(); err != nil {
		log.Fatalf("Some requests failed: %v", err)
	}

	elapsedTime := time.Since(startTime)
	log.Printf("Completed %d requests in %v", successCount, elapsedTime)
	log.Printf("QPS: %.2f", float64(successCount)/elapsedTime.Seconds())
}
