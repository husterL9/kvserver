package kvstore

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

type FileSystemAdapter struct {
	MappedFiles map[string]*MemoryMap
	store       map[string]Item
}

type MemoryMap struct {
	Data []byte
	Size int64
}

func NewFileSystemAdapter(store map[string]Item) *FileSystemAdapter {
	return &FileSystemAdapter{
		MappedFiles: make(map[string]*MemoryMap),
		store:       store,
	}
}

// 映射单个文件到内存
func (fsa *FileSystemAdapter) mapFile(path string, size int64) (*MemoryMap, error) {
	if size == 0 {
		return &MemoryMap{Data: nil, Size: 0}, nil
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := unix.Mmap(int(file.Fd()), 0, int(size), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	fsa.store[path] = Item{Key: path, Value: nil, Meta: MetaData{Type: File, Location: path}}
	return &MemoryMap{Data: data, Size: size}, nil
}

// 关闭单个内存映射
func (fsa *FileSystemAdapter) Unmap(path string, mmap *MemoryMap) {
	if mmap.Data != nil {
		if err := unix.Munmap(mmap.Data); err != nil {
			log.Printf("Failed to unmap the memory: %v", err)
		}
	}
	// 清理 Data 指针，防止后续潜在的无效内存访问
	mmap.Data = nil
	// 删除映射中的条目
	delete(fsa.MappedFiles, path)
	delete(fsa.store, path)
}

// 映射文件到内存
func (fsa *FileSystemAdapter) LoadFile(rootDir string) (map[string]*MemoryMap, error) {
	mappedFiles, err := VisitFiles(rootDir, fsa.MappedFiles, fsa.mapFile)
	fmt.Println("mappedFiles", mappedFiles)
	if err != nil {
		return mappedFiles, err
	}

	return mappedFiles, nil
}

func (fsa *FileSystemAdapter) ReadFile(path string) ([]byte, error) {
	mmap, ok := fsa.MappedFiles[path]
	if !ok {
		return nil, fmt.Errorf("file not mapped: %s", path)
	}
	if mmap.Data == nil {
		return nil, fmt.Errorf("no data mapped for file: %s", path)
	}
	// 返回一个数据副本以防止外部修改影响内存映射区域
	// dataCopy := make([]byte, mmap.Size)
	// copy(dataCopy, mmap.Data)
	return mmap.Data, nil
}

func (fsa *FileSystemAdapter) WriteFile(path string, data []byte) error {
	mmap, ok := fsa.MappedFiles[path]
	if !ok {
		return fmt.Errorf("file not mapped: %s", path)
	}
	if int64(len(data)) > mmap.Size {
		return fmt.Errorf("data size exceeds mapped size for file: %s", path)
	}
	// 更新内存映射区域
	copy(mmap.Data, data)
	return nil
}

func (fsa *FileSystemAdapter) DeleteFile(path string) error {
	mmap, ok := fsa.MappedFiles[path]
	if !ok {
		return fmt.Errorf("file not mapped: %s", path)
	}
	// 先解除映射
	err := unix.Munmap(mmap.Data)
	if err != nil {
		return fmt.Errorf("failed to unmap the memory: %v", err)
	}
	// 清理 Data 指针
	mmap.Data = nil
	// 从映射中删除条目
	delete(fsa.MappedFiles, path)
	return nil
}
