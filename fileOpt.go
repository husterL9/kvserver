package db

import (
	"fmt"
	"os"
	"path/filepath"
)

// 列出当前目录下的所有文件和文件夹
func (kv *KVStore) LsFile(currentDir string) ([]string, error) {
	var files []string
	fileList, err := os.ReadDir(currentDir)
	if err != nil {
		return nil, err // Return the error if it fails to read the directory
	}

	for _, file := range fileList {
		files = append(files, file.Name())
	}

	return files, nil
}

// 在指定目录下创建一个文件夹
func (kv *KVStore) fileMakeDir(currentDir string, dirName string) error {
	fullPath := filepath.Join(currentDir, dirName)

	err := os.Mkdir(fullPath, 0755)
	if err != nil {
		return err // Return any errors that occur
	}
	return nil
}

// 创建一个文件
func (kv *KVStore) CreateFile(path string) error {
	fmt.Println("CreateFile", path)
	// 检查文件是否已存在
	if _, err := os.Stat(path); err == nil {
		// 文件已存在，返回特定的响应或错误
		return fmt.Errorf("文件 '%s' 已存在", path)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()
	//映射到内存
	mmap, err := kv.fsAdapter.mapFile(path, 0)
	if err != nil {
		return fmt.Errorf("映射文件失败: %v", err)
	}
	kv.fsAdapter.MappedFiles[path] = mmap
	return nil
}
