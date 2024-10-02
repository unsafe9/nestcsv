package nestcsv

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
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

func isAllStringEmpty(arr []string) bool {
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

func collectStructFieldsImplementing[T any](structPtr any) []T {
	var (
		ret = make([]T, 0)
		it  = reflect.TypeOf((*T)(nil)).Elem()
		v   = reflect.ValueOf(structPtr).Elem()
	)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.IsNil() || !f.Type().Implements(it) {
			continue
		}
		ret = append(ret, f.Interface().(T))
	}
	return ret
}

func reflectObjListField(objList any, key string) (v reflect.Value, isMap bool, isMethod bool) {
	if objList == nil {
		return reflect.Value{}, false, false
	}

	v = reflect.ValueOf(objList)
	if v.Kind() != reflect.Slice {
		panic("objList is not slice")
	}
	if v.Len() == 0 {
		return v, false, false
	}
	elem := v.Index(0)
	if elem.Kind() == reflect.Interface {
		elem = elem.Elem()
	}
	if elem.Kind() == reflect.Map {
		isMap = true
	} else if m := elem.MethodByName(key); m.IsValid() {
		isMethod = true
	} else {
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		field := elem.FieldByName(key)
		if !field.IsValid() {
			panic("invalid key " + key)
		}
	}
	return v, isMap, isMethod
}

func reflectObjListFieldElemValue(v reflect.Value, idx int, isMap, isMethod bool, key string) reflect.Value {
	elem := v.Index(idx)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	if isMap {
		return elem.MapIndex(reflect.ValueOf(key))
	} else if isMethod {
		return elem.MethodByName(key).Call(nil)[0]
	} else {
		return elem.FieldByName(key)
	}
}

func sortBy(objList any, key string) any {
	v, isMap, isMethod := reflectObjListField(objList, key)
	if !v.IsValid() || v.Len() == 0 {
		return nil
	}
	sort.Slice(objList, func(i, j int) bool {
		valI := reflectObjListFieldElemValue(v, i, isMap, isMethod, key)
		valJ := reflectObjListFieldElemValue(v, j, isMap, isMethod, key)

		switch valI.Kind() {
		case reflect.String:
			return valI.String() < valJ.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return valI.Int() < valJ.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return valI.Uint() < valJ.Uint()
		case reflect.Float32, reflect.Float64:
			return valI.Float() < valJ.Float()
		default:
			return strings.Compare(valI.String(), valJ.String()) < 0
		}
	})
	return objList
}

func anyBy(objList any, key string, value any) bool {
	v, isMap, isMethod := reflectObjListField(objList, key)
	if !v.IsValid() || v.Len() == 0 {
		return false
	}
	for i := 0; i < v.Len(); i++ {
		val := reflectObjListFieldElemValue(v, i, isMap, isMethod, key)
		if val.IsValid() && val.Interface() == value {
			return true
		}
	}
	return false
}

func pascal(s string) string {
	tokens := strings.Split(strings.ReplaceAll(s, "_", " "), " ")
	for i, token := range tokens {
		if len(token) > 1 {
			tokens[i] = strings.ToUpper(token[:1]) + token[1:]
		} else if len(token) == 1 {
			tokens[i] = strings.ToUpper(token)
		}
	}
	return strings.Join(tokens, "")
}
