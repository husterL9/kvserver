package kvstore

import (
	"fmt"
	"os"
	"path/filepath"
)

// 遍历目录并对每个文件执行给定的操作
func VisitFiles(dir string, mappedFiles map[string]*MemoryMap, action func(string, int64) (*MemoryMap, error)) (map[string]*MemoryMap, error) {
	err := visit(dir, action, mappedFiles)
	return mappedFiles, err
}

// 用于遍历目录
func visit(path string, action func(string, int64) (*MemoryMap, error), mappedFiles map[string]*MemoryMap) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("错误目录")
		return err
	}
	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			err = visit(fullPath, action, mappedFiles)
			if err != nil {
				return err
			}
		} else {
			fileInfo, err := os.Stat(fullPath)
			if err != nil {
				return err
			}
			mmap, err := action(fullPath, fileInfo.Size())
			if err != nil {
				return err
			}
			mappedFiles[fullPath] = mmap
		}
	}
	return nil
}
