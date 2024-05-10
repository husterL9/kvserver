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

// Remove 从 int64 切片中删除指定的元素
func Remove(slice []int64, item int64) []int64 {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice // 如果没有找到元素，返回原切片
}
