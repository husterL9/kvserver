package kvstore

type GetArgs struct {
	Key      string
	ClientId int64 // 客户端ID
	OpId     int64 // 操作ID，用于确保操作的幂等性
}
type GetResponse struct {
	Value   []byte
	Success bool
}
