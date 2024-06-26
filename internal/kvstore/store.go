package kvstore

type Store interface {
	Get(key string) (*Item, bool)
	Set(key string, val *Item)
	// Append(key string, value []byte) (err error)
	Delete(key string)
	// Prefix(key string) map[string][]byte
	// Dump() map[string][]byte

	// Len() int
}

func GetPrefixEnd(key string) string {
	start := []byte(key)
	end := make([]byte, len(start))
	copy(end, start)
	for i := len(end) - 1; i >= 0; i-- {
		if end[i] < 0xff {
			end[i] = end[i] + 1
			end = end[:i+1]
			return string(end)
		}
	}
	return ""
}
