package db

import (
	"os"
	"testing"

	"github.com/husterL9/kvserver/internal/kvstore"
	"github.com/husterL9/kvserver/internal/wal"
)

// 创建一个简单的KVStore用于测试
func createTestKVStore() *KVStore {
	walInstance := wal.NewWAL("/tmp/test_wal")
	return &db.KVStore{
		wal: walInstance,
		// 初始化其他必要的字段
	}
}

// 清理测试环境
func cleanupTestKVStore(store *KVStore) {
	store.wal.Close()
	os.RemoveAll("/tmp/test_wal")
}
func TestGenReadView(t *testing.T) {
	store := createTestKVStore()
	defer cleanupTestKVStore(store)

	tx := &Tx{
		db:   store,
		txID: 1,
	}

	err := tx.genRV()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tx.readView.creatorTxID != tx.txID {
		t.Errorf("Expected creatorTxID to be %d, got %d", tx.txID, tx.readView.creatorTxID)
	}

	if tx.readView.lowLimitID != store.tm.nextTxID+1 {
		t.Errorf("Expected lowLimitID to be %d, got %d", store.tm.nextTxID+1, tx.readView.lowLimitID)
	}
}
func TestCommit(t *testing.T) {
	store := createTestKVStore()
	defer cleanupTestKVStore(store)

	tx := &Tx{
		db:   store,
		txID: 1,
		commits: []*Record{
			{Key: []byte("key1"), Val: []byte("value1")},
		},
	}

	err := tx.commit()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 检查activeTxIDs中是否移除了txID
	if Contains(store.tm.activeTxIDs, tx.txID) {
		t.Errorf("Expected txID %d to be removed from activeTxIDs", tx.txID)
	}

	// 检查是否正确写入WAL
	// (需要具体实现WAL的检查代码)
}
func TestRollBack(t *testing.T) {
	store := createTestKVStore()
	defer cleanupTestKVStore(store)

	tx := &Tx{
		db:   store,
		txID: 1,
		undos: []*Record{
			{Key: []byte("key1"), Val: []byte("value1")},
		},
	}

	err := tx.rollBack()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 检查是否正确加载了undo记录
	// (需要具体实现加载undo记录的检查代码)
}
func TestGet(t *testing.T) {
	store := createTestKVStore()
	defer cleanupTestKVStore(store)

	tx := &Tx{
		db:   store,
		txID: 1,
	}

	// 预设一些数据
	store.set("key1", "value1", MetaData{Type: KV_OBJ})
	store.set("key2", "value2", MetaData{Type: KV_OBJ})

	args := kvstore.GetArgs{Key: "key1"}
	response, ok := tx.Get(args)

	if !ok {
		t.Fatalf("Expected to find key, but got false")
	}

	if string(response.Value) != "value1" {
		t.Errorf("Expected value1, got %s", response.Value)
	}
}
func TestSet(t *testing.T) {
	store := createTestKVStore()
	defer cleanupTestKVStore(store)

	tx := &Tx{
		db:   store,
		txID: 1,
	}

	args := kvstore.SetArgs{Key: "key1", Value: "value1", Meta: MetaData{Type: KV_OBJ}}
	tx.Set(args)

	// 检查是否正确写入数据
	response, ok := tx.Get(kvstore.GetArgs{Key: "key1"})
	if !ok {
		t.Fatalf("Expected to find key, but got false")
	}

	if string(response.Value) != "value1" {
		t.Errorf("Expected value1, got %s", response.Value)
	}
}

//go test -v ./path/to/your/tests
