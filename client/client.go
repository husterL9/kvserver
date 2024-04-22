package client

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	pb "github.com/husterL9/kvserver/api/protobuf"
	"github.com/husterL9/kvserver/internal/kvstore"

	"google.golang.org/grpc"
)

type KVStoreClient struct {
	client   pb.KVStoreServiceClient
	clientId int64
	opId     int64 //操作id 确保操作的幂等性
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func NewKVStoreClient(conn *grpc.ClientConn) *KVStoreClient {
	return &KVStoreClient{
		client:   pb.NewKVStoreServiceClient(conn),
		clientId: nrand(),
		opId:     0,
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
	args := kvstore.GetArgs{
		Key:      key,
		ClientId: c.clientId,
		OpId:     c.opId,
	}
	c.opId++
	var value string
	var err error
	//处理重试请求
	for {
		resp, rpcErr := c.client.Get(ctx, &pb.GetRequest{Key: args.Key, ClientId: args.ClientId, OpId: args.OpId})
		err = rpcErr
		fmt.Println("resp", resp, "rpcErr", rpcErr)
		if rpcErr != nil {
			value = ""
		}
		if resp.Success {
			value = string(resp.Value)
			err = nil
			break
		}
		time.Sleep(1 * time.Second)
	}
	return value, err
}

// append
func (c *KVStoreClient) Append(key, value string, meta *pb.MetaData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.AppendRequest{
		Key:   key,
		Value: []byte(value),
		Meta:  meta,
	}

	resp, err := c.client.Append(ctx, req)
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}
