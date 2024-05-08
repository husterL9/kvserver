package kvstore

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type OpRecord struct {
	OpId   int64  // 操作ID
	OpType string // 操作类型
	Value  []byte // 操作的值
	Result []byte // 操作结果

}
type Map struct {
	kv                 map[string]Item
	fsAdapter          *FileSystemAdapter
	blockDeviceAdapter *BlockDeviceAdapter
	latestOp           map[int64]int64 // 客户端ID映射到其最新操作ID
	opHistory          map[int64][]OpRecord
}

// Item 表示存储在内存中的键值对项
type Item struct {
	Key       string
	Value     []byte
	Timestamp int64
	Committed bool
	Meta      MetaData
}

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

func NewMap() Store {
	rootDir, _ := filepath.Abs("./internal/kvstore/fakeRoot")
	kv := map[string]Item{}
	fsAdapter := NewFileSystemAdapter(kv, rootDir)
	// 异步加载文件
	go func() {
		fsAdapter.LoadFile()
	}()
	devicePath := "/dev/loop6"
	fileInfo, _ := os.Stat(devicePath)
	fmt.Println("fileInfo.Size()", fileInfo.Size())
	mmap, _ := fsAdapter.mapFile(devicePath, 100)
	fsAdapter.MappedFiles[devicePath] = mmap
	blockDeviceAdapter := NewBlockDeviceAdapter()
	return &Map{
		kv:                 kv,
		blockDeviceAdapter: blockDeviceAdapter,
		fsAdapter:          fsAdapter,
	}
}

func (m *Map) Set(key string, value []byte, meta MetaData) error {
	// m.kv[key] = val
	err := error(nil)
	switch meta.Type {
	case KVObj:
		m.kv[key] = Item{Key: key, Value: value, Meta: meta}
	case File:
		// err := kv.fsAdapter.WriteFile(meta.Location, value)
		err = error(nil)
		if err != nil {
			log.Printf("Error writing file: %v", err)
		} else {
			m.kv[key] = Item{Key: key, Value: []byte("file data at offset 0"), Meta: meta}
		}
		return err
	case BlockDevice:
		fmt.Println("meta.Offset", meta.Offset)
		err = m.blockDeviceAdapter.WriteBlock(meta.Offset, value)
		if err != nil {
			log.Printf("Error writing block device: %v", err)
		} else {
			m.kv[key] = Item{Key: key, Value: []byte("Block device data at offset 0"), Meta: meta}
		}
	default:
		log.Printf("Unsupported data type: %v", meta.Type)
	}
	return err
}

func (m *Map) Get(args GetArgs) ([]byte, bool) {
	key, clientId, opId := args.Key, args.ClientId, args.OpId
	latestOpId, ok := m.latestOp[clientId]
	item, exists := m.kv[key]
	fmt.Println("item==========", item, exists)
	// 检查是否是重复或过时的请求
	if !ok || opId > latestOpId {
		if !exists {
			return nil, false
		}
		switch item.Meta.Type {
		case KVObj:
			return item.Value, true
		case File:
			data, err := m.fsAdapter.ReadFile(key)
			if err != nil {
				log.Printf("Error reading file: %v", err)
				return nil, false
			}
			return data, true
		case BlockDevice:
			data, err := m.blockDeviceAdapter.ReadBlock(0, 3)
			if err != nil {
				log.Printf("Error reading block device: %v", err)
				return nil, false
			}
			return data, true
		}
	}
	//如果 args.OpId = latestOpId,则这个请求是重复的
	if !ok || opId == latestOpId {
		//返回最新的操作结果
		for _, record := range m.opHistory[args.ClientId] {
			if record.OpId == args.OpId {
				return record.Result, true
			}
		}
	}
	return nil, false
}
func (m *Map) Append(key string, value []byte, meta MetaData) error {

	// 检查键是否存在
	item, exists := m.kv[key]
	fmt.Println("item", item)
	if exists {
		// 检查是否为文件或块设备，因为它们的追加行为可能不同
		switch item.Meta.Type {
		case KVObj:
			// 直接追加到现有值
			item.Value = append(item.Value, value...)
			m.kv[key] = item
		case File:
			err := m.fsAdapter.AppendFile(key, value)
			if err != nil {
				return fmt.Errorf("error appending to file: %v", err)
			}
		case BlockDevice:
			// 对于块设备，你可能需要处理不同的逻辑或者不支持追加
			return fmt.Errorf("append operation not supported for block devices")
		default:
			return fmt.Errorf("unsupported data type: %v", meta.Type)
		}
		// 如果键不存在，则创建一个新的键值对
	} else {
		switch meta.Type {
		case KVObj:
			m.Set(key, value, meta)
		case File:
			// kv.store[key] = Item{Key: key, Value: nil, Meta: meta}
			//创建一个新文件并写入数据
			// err := kv.fsAdapter.WriteFile(meta.Location, value)
		case BlockDevice:
			// 对于块设备，你可能需要处理不同的逻辑或者不支持追加
			return fmt.Errorf("append operation not supported for block devices")
		default:
			return fmt.Errorf("unsupported data type: %v", meta.Type)
		}
	}
	return nil
}

func (m *Map) Delete(key string) {
	delete(m.kv, key)
}
