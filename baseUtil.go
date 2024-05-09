package db

// 判断切片中是否包含某个元素
func Contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
