package llvm

// 类型是否是
func is[T any](v any)bool{
	_, ok := v.(T)
	return ok
}

// 是否匹配
func match[T comparable](v T, to ...T)bool{
	for _, t := range to{
		if v == t{
			return true
		}
	}
	return false
}
