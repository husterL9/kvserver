// 文件路径: cmd/server/main.go

package main

import (
	"log"

	"github.com/husterL9/kvserver/internal/kvstore"
	"github.com/husterL9/kvserver/internal/server"
)

func main() {
	err := server.StartGRPCServer(kvstore.NewKVStore(), ":50051")
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
