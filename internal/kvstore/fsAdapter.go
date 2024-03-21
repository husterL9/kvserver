package kvstore

import (
	"os"
	"path/filepath"
)

type FileSystemAdapter struct {
	RootDir string
}

func NewFileSystemAdapter(rootDir string) *FileSystemAdapter {
	return &FileSystemAdapter{
		RootDir: rootDir,
	}
}

func (fsa *FileSystemAdapter) ReadFile(path string) ([]byte, error) {
	fullPath := filepath.Join(fsa.RootDir, path)
	return os.ReadFile(fullPath)
}

func (fsa *FileSystemAdapter) WriteFile(path string, data []byte) error {
	fullPath := filepath.Join(fsa.RootDir, path)
	return os.WriteFile(fullPath, data, 0644)
}

func (fsa *FileSystemAdapter) DeleteFile(path string) error {
	fullPath := filepath.Join(fsa.RootDir, path)
	return os.Remove(fullPath)
}
