// 文件路径: cmd/server/main.go

package main

import (
	"log"

	"github.com/husterL9/kvserver/internal/kvstore"
	"github.com/husterL9/kvserver/internal/server"
)

func main() {
	fsRootDir := "/home/ljw/SE8/kvserver/cmd/client/fakeRoot"                        // 文件系统根目录
	blockDevicePath := "/home/ljw/SE8/kvserver/cmd/client/fakeBlock/fakeBlockDevice" // 块设备路径
	err := server.StartGRPCServer(kvstore.NewKVStore(fsRootDir, blockDevicePath), ":50051")
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
