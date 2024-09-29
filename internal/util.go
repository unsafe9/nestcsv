package internal

import "reflect"

func ShallowCopyMap[K comparable, V any](m map[K]V) map[K]V {
	clone := make(map[K]V)
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

func FindPtr[T any](arr []*T, f func(*T) bool) *T {
	for _, v := range arr {
		if f(v) {
			return v
		}
	}
	return nil
}

func IsAllEmpty(arr []string) bool {
	for _, v := range arr {
		if v != "" {
			return false
		}
	}
	return true
}

func EqualPtr(a, b any) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}

func RemoveOne[T any](arr []T, f func(T) bool) []T {
	for i := range arr {
		if f(arr[i]) {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}
