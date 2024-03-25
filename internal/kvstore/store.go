package kvstore

import (
	"fmt"
	"log"
	"sync"
)

// DataType 表示存储在KV存储中的数据类型
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
}

// Item 表示存储在内存中的键值对项
type Item struct {
	Key   string
	Value []byte
	Meta  MetaData
}

// KVStore 表示内存键值存储引擎的主结构
type KVStore struct {
	// 使用map来存储键值对，加锁以支持并发访问
	store              map[string]Item
	lock               sync.RWMutex
	fsAdapter          *FileSystemAdapter  // 确保这一行正确无误
	blockDeviceAdapter *BlockDeviceAdapter // 如果需要处理块设备，也包括这一行
}

// NewKVStore 创建并返回一个新的KVStore实例
func NewKVStore(fsRootDir, blockDevicePath string) *KVStore {
	return &KVStore{
		store:              make(map[string]Item),
		fsAdapter:          NewFileSystemAdapter(fsRootDir),
		blockDeviceAdapter: NewBlockDeviceAdapter(blockDevicePath),
	}
}

func (kv *KVStore) Set(key string, value []byte, meta MetaData) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	fmt.Println("================", meta)
	switch meta.Type {
	case KVObj:
		kv.store[key] = Item{Key: key, Value: value, Meta: meta}
	case File:
		err := kv.fsAdapter.WriteFile(meta.Location, value)
		if err != nil {
			log.Printf("Error writing file: %v", err)
		} else {
			kv.store[key] = Item{Key: key, Value: value, Meta: meta}
		}
	case BlockDevice:
		// 为简化示例，这里直接将偏移量设置为0，实际应用中可能需要更复杂的逻辑
		err := kv.blockDeviceAdapter.WriteBlock(0, value)
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

func (kv *KVStore) Get(key string) (Item, bool) {
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	item, exists := kv.store[key]
	if !exists {
		return Item{}, false
	}

	switch item.Meta.Type {
	case File:
		data, err := kv.fsAdapter.ReadFile(item.Meta.Location)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			return Item{}, false
		}
		item.Value = data
	case BlockDevice:
		// 为简化示例，这里直接将偏移量设置为0，并假定读取长度为1024，实际应用中可能需要更复杂的逻辑
		data, err := kv.blockDeviceAdapter.ReadBlock(0, 1024)
		if err != nil {
			log.Printf("Error reading block device: %v", err)
			return Item{}, false
		}
		item.Value = data
	}

	return item, true
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
		// 对于普通键值对，直接从存储中删除
		delete(kv.store, key)
	case File:
		// 对于文件类型，需要删除文件系统中的文件
		err := kv.fsAdapter.DeleteFile(item.Meta.Location)
		if err != nil {
			return fmt.Errorf("error deleting file: %v", err)
		}
		// 从存储中删除条目
		delete(kv.store, key)
	case BlockDevice:
		// 对于块设备，可能需要执行特定的清除操作
		// 注意：这里我们只是示意性地展示了接口，实际上块设备不应简单地“删除”
		// 可能需要根据实际场景实现适当的逻辑，这里不做具体实现
		log.Printf("Block device delete operation is not implemented")
		// 仍然从存储中删除条目，即使块设备的“删除”可能并未实际执行
		delete(kv.store, key)
	default:
		return fmt.Errorf("unsupported data type: %v", item.Meta.Type)
	}

	return nil
}
