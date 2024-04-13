package kvstore_test

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/husterL9/kvserver/internal/kvstore"
	"github.com/stretchr/testify/require"
)

// 测试Set和Get操作
func TestSetGet(t *testing.T) {
	kv := kvstore.NewKVStore() // 初始化你的KV存储实例
	key := "testKey"
	expectedValue := []byte("testValue")

	// 测试Set操作
	err := kv.Set(key, expectedValue, kvstore.MetaData{
		Type:     kvstore.KVObj,
		Location: "",
		Offset:   0,
		Size:     0,
	})
	require.NoError(t, err)

	// 测试Get操作
	value, err := kv.Get(key)
	require.NoError(t, err)
	require.Equal(t, expectedValue, value)

	// 清理
	err = kv.Delete(key)
	require.NoError(t, err)
}

// 测试Delete操作
func TestDelete(t *testing.T) {
	kv := NewKVStore() // 初始化你的KV存储实例
	key := "testKey"
	value := []byte("testValue")

	// 先设置一个键值对
	err := kv.Set(key, value)
	require.NoError(t, err)

	// 删除键值对
	err = kv.Delete(key)
	require.NoError(t, err)

	// 尝试获取已删除的键值对
	_, err = kv.Get(key)
	require.Error(t, err)
}

func TestFileAccess(t *testing.T) {
	filePath := "/tmp/testfile"
	content := []byte("Hello, file!")

	// 写入文件
	err := os.WriteFile(filePath, content, 0644)
	require.NoError(t, err)

	// 读取并验证内容
	readContent, err := ioutil.ReadFile(filePath)
	require.NoError(t, err)
	require.Equal(t, content, readContent)

	// 清理文件
	os.Remove(filePath)
}

// 块设备访问测试类似，但需要具有块设备路径和适当的权限
//
// 混合操作测试
// 结合使用Set、Get和Delete操作，并可能涉及到并发访问
func TestMixedOperations(t *testing.T) {
	kv := NewKVStore() // 初始化你的KV存储实例
	keys := []string{"key1", "key2", "key3"}
	values := [][]byte{[]byte("value1"), []byte("value2"), []byte("value3")}

	var wg sync.WaitGroup

	// 并发设置键值对
	for i, key := range keys {
		wg.Add(1)
		go func(key string, value []byte) {
			defer wg.Done()
			err := kv.Set(key, value)
			require.NoError(t, err)
		}(key, values[i])
	}

	wg.Wait()

	// 并发读取并验证键值对
	for i, key := range keys {
		wg.Add(1)
		go func(key string, expectedValue []byte) {
			defer wg.Done()
			value, err := kv.Get(key)
			require.NoError(t, err)
			require.Equal(t, expectedValue, value)
		}(key, values[i])
	}

	wg.Wait()

	// 并发删除键值对
	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			err := kv.Delete(key)
			require.NoError(t, err)
		}(key)
	}

	wg.Wait()
}
