package db

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/husterL9/kvserver/api/protobuf"
	"github.com/husterL9/kvserver/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestPerformance(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := client.NewKVStoreClient(conn)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	// defer c.Close()

	// 测试写入性能
	start := time.Now()
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", rand.Intn(1000000))
		ok, err := c.Set(key, value, &protobuf.MetaData{})
		if err != nil {
			log.Printf("Failed to set key: %v", err)
			log.Printf("ok: %v", ok)
		}

	}
	elapsed := time.Since(start)
	log.Printf("Write 10000 records took %s", elapsed)

	// 测试读取性能
	start = time.Now()
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", rand.Intn(10000))
		_, err := c.Get(key)
		if err != nil {
			log.Printf("Failed to get key: %v", err)
		}
	}
	elapsed = time.Since(start)
	log.Printf("Read 10000 records took %s", elapsed)
}
