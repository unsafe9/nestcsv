package internal

func ShallowCopyMap[K comparable, V any](m map[K]V) map[K]V {
	clone := make(map[K]V)
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

func Some[T any](arr []T, f func(T) bool) bool {
	for _, v := range arr {
		if f(v) {
			return true
		}
	}
	return false
}

func FindPtr[T any](arr []*T, f func(*T) bool) *T {
	for _, v := range arr {
		if f(v) {
			return v
		}
	}
	return nil
}
