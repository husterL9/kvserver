package db

import (
	"sync"

	"github.com/husterL9/kvserver/internal/kvstore"
)

// KVStore 表示内存键值存储引擎的主结构
type KVStore struct {
	// 使用map来存储键值对，加锁以支持并发访问
	store kvstore.Store
	lock  sync.RWMutex
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
		c.Path = "moss"
	}
	if c.Store == nil {
		c.Store = kvstore.NewMap()
	}
	return c
}

// NewKVStore 创建并返回一个新的KVStore实例
func NewKVStore(conf *Config) *KVStore {
	conf = conf.Default()
	store := conf.Store
	db := &KVStore{
		store: store,
	}
	return db
}

func (db *KVStore) Set(key string, value []byte, meta kvstore.MetaData) {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.store.Set(key, value, meta)

}

func (db *KVStore) Get(args kvstore.GetArgs) (kvstore.GetResponse, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	val, ok := db.store.Get(args)
	return kvstore.GetResponse{
		Value: val,
	}, ok
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
