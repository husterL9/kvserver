package kvstore

type GetArgs struct {
	Key      string
	ClientId int64 // 客户端ID
	OpId     int64 // 操作ID，确保幂等性
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
