package client

import (
	"context"
	"time"

	pb "github.com/husterL9/kvserver/api/protobuf"
)

func (c *KVStoreClient) LsFile(currentDir string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := c.client.FileLs(ctx, &pb.FileLsRequest{
		Path: currentDir,
	})
	if err != nil {
		return nil, err
	}
	return resp.Files, nil
}
