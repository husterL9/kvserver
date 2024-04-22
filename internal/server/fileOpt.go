package server

import (
	"context"
	"os"

	pb "github.com/husterL9/kvserver/api/protobuf"
)

// LsFile
func (s *server) FileLs(ctx context.Context, req *pb.FileLsRequest) (*pb.FileLsResponse, error) {
	path := req.GetPath()
	files, err := s.store.LsFile(path)
	if err != nil {
		return nil, err // Return the error to the client
	}

	resp := &pb.FileLsResponse{
		Files: files,
	}
	return resp, nil
}

// CdDir
func (s *server) CdDir(ctx context.Context, req *pb.CdDirRequest) (*pb.CdDirResponse, error) {
	dir := req.GetPath()
	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, err // 如果目录不存在，返回错误
	}

	// Change the server's current directory
	// s.currentDir = newDir
	return &pb.CdDirResponse{
		Success: true,
	}, nil
}

// MakeDir
func (s *server) MakeDir(ctx context.Context, req *pb.MakeDirRequest) (*pb.MakeDirResponse, error) {
	dir := req.GetPath()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &pb.MakeDirResponse{
		Success: true,
	}, nil
}
