package server

import (
	"context"
	"log"
	"net"

	pb "github.com/husterL9/kvserver/api/protobuf" // 替换为你的protobuf包路径

	"github.com/husterL9/kvserver/internal/kvstore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server是KVStoreService的实现
type server struct {
	store *kvstore.KVStore
	pb.UnimplementedKVStoreServiceServer
}

// NewServer
func NewServer(store *kvstore.KVStore) *server {
	return &server{store: store}
}

// Set
func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	meta := kvstore.MetaData{
		Type:     kvstore.DataType(req.Meta.Type),
		Location: req.Meta.Location,
		Offset:   req.Meta.Offset,
		Size:     req.Meta.Size,
	}
	s.store.Set(req.Key, req.Value, meta)
	return &pb.SetResponse{Success: true}, nil
}

// Get
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	args := kvstore.GetArgs{
		Key:      req.GetKey(),
		ClientId: req.GetClientId(),
		OpId:     req.GetOpId(),
	}
	// 从KVStore中检索键对应的值
	item, exists := s.store.Get(args)
	if !exists {
		// 如果键不存在，可以返回一个错误或一个空的响应
		log.Printf("Key not found: %s", req.GetKey())
		return &pb.GetResponse{Value: nil, Success: exists}, status.Errorf(codes.NotFound, "key not found: %s", req.GetKey())
	}
	// 如果键存在，返回找到的值
	return &pb.GetResponse{Value: item.Value, Success: exists}, nil
}

// append
// func (s *server) Append(ctx context.Context, req *pb.AppendRequest) (*pb.AppendResponse, error) {
// 	args := kvstore.AppendArgs{
// 		Key:   req.GetKey(),
// 		Value: req.GetValue(),
// 	}
// }

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
