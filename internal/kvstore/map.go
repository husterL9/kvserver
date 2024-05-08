package kvstore

import (
	"fmt"
)

type Map struct {
	kv map[string]*Item
}

// Item 表示存储在内存中的键值对项
type Item struct {
	Key       string
	Value     []byte
	Timestamp int64
	Committed bool
	Meta      MetaData
	Next      *Item
}

func NewMap() Store {

	kv := map[string]*Item{}

	return &Map{
		kv: kv,
	}
}

func (m *Map) Set(key string, value []byte) error {
	// m.kv[key] = val

}

func (m *Map) Get(key string) (*Item, bool) {
	val, ok := m.kv[key]
	if !ok {
		return nil, false
	}
	return val, true
}

func (m *Map) Append(key string, value []byte) error {

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
