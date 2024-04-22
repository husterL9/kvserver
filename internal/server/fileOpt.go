package server

import (
	"context"

	pb "github.com/husterL9/kvserver/api/protobuf"
)

// LsFile
func (s *server) LsFile(ctx context.Context, req *pb.FileLsRequest) (*pb.FileLsResponse, error) {
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
