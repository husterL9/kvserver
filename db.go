package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/husterL9/kvserver/internal/kvstore"
	"github.com/husterL9/kvserver/internal/wal"
)

type OpRecord struct {
	OpId   int64  // 操作ID
	OpType string // 操作类型
	Value  []byte // 操作的值
	Result []byte // 操作结果

}

// KVStore 表示内存键值存储引擎的主结构
type KVStore struct {
	lock sync.RWMutex
	// 使用map来存储键值对，加锁以支持并发访问
	store              kvstore.Store
	tm                 *TxManager
	wal                *wal.Log
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
	fsAdapter := NewFileSystemAdapter(store, rootDir)
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

func (db *KVStore) Set(key string, value []byte, meta kvstore.MetaData) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	fmt.Println("Set key:", key, "value:", string(value), "meta:", meta)
	err := db.set(key, value, meta)
	return err
}

func (db *KVStore) set(key string, value []byte, meta kvstore.MetaData) error {
	err := error(nil)
	version := &kvstore.Version{
		Value: value,
		Meta:  meta,
	}
	item, ok := db.store.Get(key)
	if ok {
		//插入新版本
		oldVersion := item.Version
		version.Next = oldVersion
		item.Version = version
		return nil
	} else {
		fmt.Println("meta.Type", meta.Type)
		switch meta.Type {
		case kvstore.KVObj:
			db.store.Set(key, &kvstore.Item{Key: key, Version: version})
		case kvstore.File:
			err := db.fsAdapter.WriteFile(meta.Location, value)
			if err != nil {
				log.Printf("Error writing file: %v", err)
			} else {
				db.store.Set(key, &kvstore.Item{Key: key, Version: &kvstore.Version{
					Value: []byte("file data at offset 0"),
					Meta:  meta,
				}})
			}
			return err
		case kvstore.BlockDevice:
			fmt.Println("meta.Offset", meta.Offset)
			err = db.blockDeviceAdapter.WriteBlock(meta.Offset, value)
			if err != nil {
				log.Printf("Error writing block device: %v", err)
			} else {
				db.store.Set(key, &kvstore.Item{Key: key, Version: &kvstore.Version{
					Value: []byte("Block device data at offset 0"),
					Meta:  meta,
				}})
			}
		default:
			log.Printf("Unsupported data type: %v", meta.Type)
		}
		return err
	}

}
func (db *KVStore) Get(args kvstore.GetArgs) (kvstore.GetResponse, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	key, clientId, opId := args.Key, args.ClientId, args.OpId
	latestOpId, ok := db.latestOp[clientId]
	item, exists := db.store.Get(key)
	firstVersion := item.Version
	fmt.Println("item==========", item, exists)
	if !exists {
		return kvstore.GetResponse{
			Value: nil,
		}, false
	}
	// 检查是否是重复或过时的请求
	if !ok || opId > latestOpId {
		fmt.Println("firstVersion.Meta.Type", firstVersion)
		switch firstVersion.Meta.Type {
		case kvstore.KVObj:
			return kvstore.GetResponse{
				Value: firstVersion.Value,
			}, true
		case kvstore.File:
			data, err := db.fsAdapter.ReadFile(key)
			if err != nil {
				log.Printf("Error reading file: %v", err)
				return kvstore.GetResponse{
					Value: nil,
				}, false
			}
			return kvstore.GetResponse{
				Value: data,
			}, true
		case kvstore.BlockDevice:
			data, err := db.blockDeviceAdapter.ReadBlock(0, 3)
			if err != nil {
				log.Printf("Error reading block device: %v", err)
				return kvstore.GetResponse{
					Value: nil,
				}, false
			}
			return kvstore.GetResponse{
				Value: data,
			}, true
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

func (db *KVStore) get(args kvstore.GetArgs) (*kvstore.Item, bool) {
	item, ok := db.store.Get(args.Key)
	return item, ok
}
func (db *KVStore) Append(key string, value []byte, meta kvstore.MetaData) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	// 检查键是否存在
	item, exists := db.store.Get(key)
	version := &kvstore.Version{Meta: meta}
	if exists {
		firstVersion := item.Version
		newVal := append(firstVersion.Value, value...)
		db.set(key, newVal, meta)
	} else {
		version.Value = value
		db.set(key, value, meta)
	}
	return nil
}

// Delete 根据键删除一个键值对
func (db *KVStore) Delete(key string) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.store.Delete(key)
	return nil
}
