// 文件路径: cmd/server/main.go

package main

import (
	"log"

	db "github.com/husterL9/kvserver"
	"github.com/husterL9/kvserver/server"
)

func main() {
	// fsRootDir := ""       // 文件系统根目录
	// blockDevicePath := "" // 块设备路径

	config := &db.Config{}
	err := server.StartGRPCServer(db.NewKVStore(config), ":50051")
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
