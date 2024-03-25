// 文件路径: client/client.go

package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/husterL9/kvserver/internal/api/protobuf"

	"google.golang.org/grpc"
)

type KVStoreClient struct {
	client pb.KVStoreServiceClient
}

func NewKVStoreClient(conn *grpc.ClientConn) *KVStoreClient {
	return &KVStoreClient{
		client: pb.NewKVStoreServiceClient(conn),
	}
}

func (c *KVStoreClient) Set(key, value string, meta *pb.MetaData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.Set(ctx, &pb.SetRequest{
		Key:   key,
		Value: []byte(value),
		Meta:  meta,
	})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (c *KVStoreClient) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.Get(ctx, &pb.GetRequest{Key: key})
	fmt.Println("resp===", resp)
	if err != nil {
		return "", err
	}
	return string(resp.Value), nil
}
