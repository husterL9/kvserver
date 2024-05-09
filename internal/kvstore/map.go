package kvstore

type Map struct {
	kv map[string]*Item
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

// Item 表示存储在内存中的键值对项
type Item struct {
	Key       string
	Value     []byte
	TxID      int64
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

func (m *Map) Set(key string, val *Item) {
	m.kv[key] = val
}

func (m *Map) Get(key string) (*Item, bool) {
	val, ok := m.kv[key]
	if !ok {
		return nil, false
	}
	return val, true
}

// func (m *Map) Append(key string, value []byte) (err error) {
// 	item, ok := m.Get(key)
// 	if ok {
// 		item.Value = append(item.Value, value...)
// 		m.kv[key] = item
// 	}
// }

func (m *Map) Delete(key string) {
	delete(m.kv, key)
}
