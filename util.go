package nestcsv

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

func mapWithoutKey[K comparable, V any](m map[K]V, key K) map[K]V {
	clone := make(map[K]V)
	for k, v := range m {
		if k != key {
			clone[k] = v
		}
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

func createFile(rootDir, fileName, ext string) (*os.File, error) {
	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create the directory: %s, %w", rootDir, err)
	}

	ext = "." + strings.TrimPrefix(ext, ".")
	fileName = strings.TrimSuffix(fileName, ext) + ext

	file, err := os.Create(filepath.Join(rootDir, fileName))
	if err != nil {
		return nil, fmt.Errorf("failed to create the file: %s, %w", fileName, err)
	}
	return file, nil
}

func saveCSVFile(rootDir, fileName string, csvData [][]string) error {
	file, err := createFile(rootDir, fileName, "csv")
	if err != nil {
		return fmt.Errorf("failed to create the file: %s, %w", fileName, err)
	}
	defer file.Close()

	if err := csv.NewWriter(file).WriteAll(csvData); err != nil {
		return fmt.Errorf("failed to write the file: %s, %w", fileName, err)
	}
	return nil
}

func saveJSONFile(rootDir, fileName string, jsonBytes []byte) error {
	file, err := createFile(rootDir, fileName, "json")
	if err != nil {
		return fmt.Errorf("failed to create the file: %s, %w", fileName, err)
	}
	defer file.Close()

	if _, err := file.Write(jsonBytes); err != nil {
		return fmt.Errorf("failed to write the file: %s, %w", fileName, err)
	}
	return nil
}

func saveBinFile(rootDir, fileName string, jsonBytes []byte) error {
	file, err := createFile(rootDir, fileName, "bin")
	if err != nil {
		return fmt.Errorf("failed to create the file: %s, %w", fileName, err)
	}
	defer file.Close()

	binHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(binHeader, uint32(len(jsonBytes)))
	if _, err := file.Write(binHeader); err != nil {
		return fmt.Errorf("failed to write the binary header: %s, %w", fileName, err)
	}

	if _, err := file.Write(jsonBytes); err != nil {
		return fmt.Errorf("failed to write the file: %s, %w", fileName, err)
	}
	return nil
}
