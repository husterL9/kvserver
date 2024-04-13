package kvstore

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
)

// 存储在KV存储中的数据类型
type DataType int

const (
	// KVObj 表示键值对对象
	KVObj DataType = iota
	// File 表示文件对象
	File
	// BlockDevice 表示块设备对象
	BlockDevice
)

// MetaData 代表存储在KV存储中的每个对象的元数据
type MetaData struct {
	Type     DataType // 对象类型: KVObj, File, 或 BlockDevice
	Location string   // 文件或块设备的位置信息（对于KV对象，此字段可以为空）
	Offset   int64    // 读/写操作的起始偏移量
	Size     int64    // 读/写操作的数据长度
}

// Item 表示存储在内存中的键值对项
type Item struct {
	Key   string
	Value []byte
	Meta  MetaData
}
type OpRecord struct {
	OpId   int64  // 操作ID
	OpType string // 操作类型
	Value  string // 操作的值
	Result string // 操作结果
}

// KVStore 表示内存键值存储引擎的主结构
type KVStore struct {
	// 使用map来存储键值对，加锁以支持并发访问
	store              map[string]Item
	lock               sync.RWMutex
	fsAdapter          *FileSystemAdapter  // 确保这一行正确无误
	blockDeviceAdapter *BlockDeviceAdapter // 如果需要处理块设备，也包括这一行
	latestOp           map[int64]int64     // 客户端ID映射到其最新操作ID
	opHistory          map[int64][]OpRecord
}

// NewKVStore 创建并返回一个新的KVStore实例
func NewKVStore() *KVStore {
	rootDir, _ := filepath.Abs("./internal/kvstore/fakeRoot")
	store := make(map[string]Item)
	fsAdapter := NewFileSystemAdapter(store)
	fsAdapter.LoadFile(rootDir)
	blockDeviceAdapter := NewBlockDeviceAdapter()
	kv := &KVStore{
		store:              store,
		fsAdapter:          fsAdapter,
		blockDeviceAdapter: blockDeviceAdapter,
	}
	return kv
}

func (kv *KVStore) Set(key string, value []byte, meta MetaData) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	switch meta.Type {
	case KVObj:
		kv.store[key] = Item{Key: key, Value: value, Meta: meta}
	case File:
		// err := kv.fsAdapter.WriteFile(meta.Location, value)
		err := error(nil)
		if err != nil {
			log.Printf("Error writing file: %v", err)
		} else {
			kv.store[key] = Item{Key: key, Value: []byte("file data at offset 0"), Meta: meta}
		}
	case BlockDevice:
		fmt.Println("meta.Offset", meta.Offset)
		err := kv.blockDeviceAdapter.WriteBlock(meta.Offset, value)
		if err != nil {
			log.Printf("Error writing block device: %v", err)
		} else {
			// 块设备写入成功，保存引用和元数据到KV存储中
			// 注意：这里我们可能会存储一个表示数据位置或描述的信息，而不是数据本身
			// 这是因为直接从块设备读取数据可能需要特定的上下文或操作
			kv.store[key] = Item{Key: key, Value: []byte("Block device data at offset 0"), Meta: meta}
		}
	default:
		log.Printf("Unsupported data type: %v", meta.Type)
	}
}

func (kv *KVStore) Get(args GetArgs) (GetResponse, bool) {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	key, clientId, opId := args
	item, exists := kv.store[key]
	if !exists {
		return GetResponse{}, false
	}
	resp := GetResponse{}
	switch item.Meta.Type {
	case File:
		data, err := kv.fsAdapter.ReadFile(key)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			return GetResponse{}, false
		}
		resp.Value = data
	case BlockDevice:
		data, err := kv.blockDeviceAdapter.ReadBlock(0, 3)
		if err != nil {
			log.Printf("Error reading block device: %v", err)
			return GetResponse{}, false
		}
		resp.Value = data
	}

	return resp, true
}

// Delete 根据键删除一个键值对
func (kv *KVStore) Delete(key string) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	item, exists := kv.store[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	switch item.Meta.Type {
	case KVObj:
		delete(kv.store, key)
	case File:
		// err := kv.fsAdapter.DeleteFile(item.Meta.Location)
		err := error(nil)
		if err != nil {
			return fmt.Errorf("error deleting file: %v", err)
		}
		delete(kv.store, key)
	case BlockDevice:
		log.Printf("Block device delete operation is not implemented")
		delete(kv.store, key)
	default:
		return fmt.Errorf("unsupported data type: %v", item.Meta.Type)
	}

	return nil
}
