package kvstore

type GetArgs struct {
	Key      string
	ClientId int64 // 客户端ID
	OpId     int64 // 操作ID，确保幂等性
}
type SetArgs struct {
	Key      string
	Value    []byte
	ClientId int64
	OpId     int64
	Meta     MetaData
}
type GetResponse struct {
	Value   []byte
	Success bool
}

type AppendArgs struct {
	Key   string
	Value []byte
}

type AppendResponse struct {
	Success bool
}
