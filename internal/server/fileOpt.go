package server

import (
	"context"
	"fmt"
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

	// 检查目录是否已存在
	if _, err := os.Stat(dir); err == nil {
		// 目录已存在，返回特定的响应或错误
		return nil, fmt.Errorf("目录 '%s' 已存在", dir)
	} else if !os.IsNotExist(err) {
		// 出现了非"不存在"的其它错误，直接返回这个错误
		return nil, err
	}

	// 尝试创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		// 创建目录失败，返回错误
		return nil, err
	}

	// 目录成功创建
	return &pb.MakeDirResponse{
		Success: true,
	}, nil
}
