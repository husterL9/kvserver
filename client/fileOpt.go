package client

import (
	"context"
	"time"

	pb "github.com/husterL9/kvserver/api/protobuf"
)

// LsFile
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

// CdDir
func (c *KVStoreClient) CdDir(dir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.client.CdDir(ctx, &pb.CdDirRequest{
		Path: dir,
	})
	return err
}

// MakeDir
func (c *KVStoreClient) MakeDir(dir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.client.MakeDir(ctx, &pb.MakeDirRequest{
		Path: dir,
	})
	return err
}
