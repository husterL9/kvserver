// 文件路径: cmd/client/main.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/husterL9/kvserver/internal/client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: client [key] [value]", os.Args)
		os.Exit(1)
	}
	key := os.Args[1]
	value := os.Args[2]

	// 连接到gRPC服务
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := client.NewKVStoreClient(conn)

	// 设置键值对
	success, err := c.Set(key, value)
	if err != nil {
		log.Fatalf("could not set key-value: %v", err)
	}
	fmt.Printf("Set result: %v\n", success)

	// 获取键值对
	gotValue, err := c.Get(key)
	if err != nil {
		log.Fatalf("could not get value: %v", err)
	}
	fmt.Printf("Got value: %s\n", gotValue)
}
