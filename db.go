package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/husterL9/kvserver/internal/kvstore"
)

type OpRecord struct {
	OpId   int64  // 操作ID
	OpType string // 操作类型
	Value  []byte // 操作的值
	Result []byte // 操作结果

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

// KVStore 表示内存键值存储引擎的主结构
type KVStore struct {
	lock sync.RWMutex
	// 使用map来存储键值对，加锁以支持并发访问
	store kvstore.Store
	tm    *TxManager

	latestOp           map[int64]int64 // 客户端ID映射到其最新操作ID
	opHistory          map[int64][]OpRecord
	blockDeviceAdapter *BlockDeviceAdapter
	fsAdapter          *FileSystemAdapter
}
type Config struct {
	Path  string
	Store kvstore.Store
}

func (c *Config) Default() *Config {
	if c == nil {
		c = &Config{}
	}
	if c.Path == "" {
		c.Path = "db"
	}
	if c.Store == nil {
		c.Store = kvstore.NewMap()
	}
	return c
}

// NewKVStore 创建并返回一个新的KVStore实例
func NewKVStore(conf *Config) *KVStore {
	rootDir, _ := filepath.Abs("./internal/kvstore/fakeRoot")
	conf = conf.Default()
	store := conf.Store
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
	db := &KVStore{
		store:              store,
		tm:                 NewTxManager(),
		blockDeviceAdapter: blockDeviceAdapter,
		fsAdapter:          fsAdapter,
	}
	return db
}

func (db *KVStore) Set(key string, value []byte, meta MetaData) {
	db.lock.Lock()
	defer db.lock.Unlock()
	err := error(nil)
	switch meta.Type {
	case KVObj:
		db.kv[key] = &Item{Key: key, Value: value, Meta: meta}
	case File:
		// err := kv.fsAdapter.WriteFile(meta.Location, value)
		err = error(nil)
		if err != nil {
			log.Printf("Error writing file: %v", err)
		} else {
			db.kv[key] = &Item{Key: key, Value: []byte("file data at offset 0"), Meta: meta}
		}
		return err
	case BlockDevice:
		fmt.Println("meta.Offset", meta.Offset)
		err = m.blockDeviceAdapter.WriteBlock(meta.Offset, value)
		if err != nil {
			log.Printf("Error writing block device: %v", err)
		} else {
			m.kv[key] = &Item{Key: key, Value: []byte("Block device data at offset 0"), Meta: meta}
		}
	default:
		log.Printf("Unsupported data type: %v", meta.Type)
	}
	return err
	db.store.Set(key, value, meta)

}

func (db *KVStore) Get(args kvstore.GetArgs) (kvstore.GetResponse, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	key, clientId, opId := args.Key, args.ClientId, args.OpId
	latestOpId, ok := db.latestOp[clientId]
	item, exists := db.store.Get(key)
	fmt.Println("item==========", item, exists)
	// 检查是否是重复或过时的请求
	if !ok || opId > latestOpId {
		if !exists {
			return kvstore.GetResponse{
				Value: nil,
			}, false
		}
		switch item.Meta.Type {
		case kvstore.KVObj:
			return kvstore.GetResponse{
				Value: item.Value,
			}, true
		case kvstore.File:
			data, err := db.store.fsAdapter.ReadFile(key)
			if err != nil {
				log.Printf("Error reading file: %v", err)
				return nil, false
			}
			return data, true
		case kvstore.BlockDevice:
			data, err := db.blockDeviceAdapter.ReadBlock(0, 3)
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
		for _, record := range db.opHistory[args.ClientId] {
			if record.OpId == args.OpId {
				return kvstore.GetResponse{
					Value: record.Result}, true
			}
		}
	}
	return kvstore.GetResponse{
		Value: nil,
	}, false

}

func (db *KVStore) Append(key string, value []byte, meta kvstore.MetaData) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	err := db.store.Append(key, value, meta)
	return err
}

// Delete 根据键删除一个键值对
func (db *KVStore) Delete(key string) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.store.Delete(key)
	return nil
}
