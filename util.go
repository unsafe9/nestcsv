package nestcsv

import (
	"io/fs"
	"iter"
	"path/filepath"
	"reflect"
	"slices"
)

func shallowCopyMap[K comparable, V any](m map[K]V) map[K]V {
	clone := make(map[K]V)
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

func findPtr[T any](arr []*T, f func(*T) bool) *T {
	for _, v := range arr {
		if f(v) {
			return v
		}
	}
	return nil
}

func isAllEmpty(arr []string) bool {
	for _, v := range arr {
		if v != "" {
			return false
		}
	}
	return true
}

func equalPtr(a, b any) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}

func removeOne[T any](arr []T, f func(T) bool) []T {
	for i := range arr {
		if f(arr[i]) {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

func walkFiles(dirs []string, files []string, exts []string) iter.Seq[string] {
	return func(yield func(string) bool) {
		visited := make(map[string]struct{})
		stop := false

		for _, dir := range dirs {
			err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				if len(exts) > 0 {
					ext := filepath.Ext(path)
					if len(ext) > 0 && !slices.Contains(exts, ext[1:]) {
						return nil
					}
				}

				path, err = filepath.Rel(".", path)
				if err != nil {
					return nil
				}

				if _, ok := visited[path]; ok {
					return nil
				}
				visited[path] = struct{}{}

				if !yield(path) {
					stop = true
					return filepath.SkipAll
				}

				return nil
			})
			if err != nil {
				return
			}
			if stop {
				return
			}
		}

		for _, path := range files {
			var err error
			path, err = filepath.Rel(".", path)
			if err != nil {
				continue
			}

			if _, ok := visited[path]; ok {
				continue
			}
			visited[path] = struct{}{}

			if !yield(path) {
				return
			}
		}
	}
}
