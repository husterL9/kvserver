package db

import (
	"fmt"
	"log"
	"os"

	"github.com/husterL9/kvserver/internal/kvstore"
	"golang.org/x/sys/unix"
)

type FileSystemAdapter struct {
	rootDir     string
	MappedFiles map[string]*MemoryMap
	store       kvstore.Store
}

type MemoryMap struct {
	Data []byte
	Size int64
}

func NewFileSystemAdapter(store kvstore.Store, rootDir string) *FileSystemAdapter {
	return &FileSystemAdapter{
		MappedFiles: make(map[string]*MemoryMap),
		store:       store,
		rootDir:     rootDir,
	}
}

// 映射单个文件到内存
func (fsa *FileSystemAdapter) mapFile(path string, size int64) (*MemoryMap, error) {
	if path == "/home/ljw/SE8/kvserver/internal/kvstore/fakeBlock/fakeBlockDevice" {
		fmt.Println("path", path)
	}
	version := &kvstore.Version{
		Meta: kvstore.MetaData{Type: kvstore.File, Location: path},
	}
	if size == 0 {
		fsa.store.Set(path, &kvstore.Item{Key: path, Version: version})
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
	fsa.store.Set(path, &kvstore.Item{Key: path, Version: version})
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
	fsa.store.Delete(path)
}

// 映射文件到内存
func (fsa *FileSystemAdapter) LoadFile() (map[string]*MemoryMap, error) {
	mappedFiles, err := VisitFiles(fsa.rootDir, fsa.MappedFiles, fsa.mapFile)
	if err != nil {
		return mappedFiles, err
	}
	return mappedFiles, nil
}

func (fsa *FileSystemAdapter) ReadFile(path string) ([]byte, error) {
	mmap, ok := fsa.MappedFiles[path]
	fmt.Println("path", path)
	if path == "/dev/loop6" {
		fmt.Println("path", path)
		return mmap.Data[:5], nil
	}

	if !ok {
		return nil, fmt.Errorf("file not mapped: %s", path)
	}
	if mmap.Data == nil {
		return []byte(""), nil
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
	copy(mmap.Data[6:10], data)
	return nil
}

func (fsa *FileSystemAdapter) AppendFile(path string, data []byte) error {
	if path == "/dev/loop6" {
		fmt.Println("path", path)
		fsa.WriteFile(path, data)
		return nil
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取当前文件大小
	fileInfo, err := file.Stat()
	fmt.Println("fileInfo", fileInfo.Size())
	if err != nil {
		return err
	}
	currentSize := fileInfo.Size()

	// 计算新的文件大小
	newSize := currentSize + int64(len(data))

	// 扩展文件大小
	if err := file.Truncate(newSize); err != nil {
		return err
	}

	// 重新映射文件
	mmap, mapErr := fsa.mapFile(path, newSize)
	if mapErr != nil {
		return err
	}
	fsa.MappedFiles[path] = mmap

	// 将数据追加到映射区域
	copy(mmap.Data[currentSize:], data)

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
