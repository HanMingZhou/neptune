package math

// Min 返回两个 int64 的最小值
func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Min3 返回三个 int64 的最小值
func Min3(a, b, c int64) int64 {
	return Min(Min(a, b), c)
}
