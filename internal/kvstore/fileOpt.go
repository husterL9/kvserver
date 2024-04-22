package kvstore

import (
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

// 创建当前目录下一个新文件
// func CreateFile(currentDir string) error {

// }
