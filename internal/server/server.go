package server

import (
	"context"
	"log"
	"net"

	pb "github.com/husterL9/kvserver/internal/api/protobuf" // 替换为你的protobuf包路径

	"github.com/husterL9/kvserver/internal/kvstore"

	"google.golang.org/grpc"
)

// server是KVStoreService的实现
type server struct {
	store *kvstore.KVStore
	pb.UnimplementedKVStoreServiceServer
}

// NewServer创建一个gRPC服务的实例
func NewServer(store *kvstore.KVStore) *server {
	return &server{store: store}
}

// Set实现了KVStoreService的Set方法
func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	meta := kvstore.MetaData{
		Type:     kvstore.DataType(req.Meta.Type),
		Location: req.Meta.Location,
	}
	s.store.Set(req.Key, req.Value, meta)
	return &pb.SetResponse{Success: true}, nil
}

// 启动gRPC服务器
func StartGRPCServer(store *kvstore.KVStore, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterKVStoreServiceServer(s, NewServer(store))
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}
