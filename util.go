package nestcsv

import (
	"embed"
	"encoding/csv"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/gertd/go-pluralize"
	"iter"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"text/template"
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

func appendUnique[T comparable](arr []T, v ...T) []T {
	for _, vv := range v {
		if !slices.Contains(arr, vv) {
			arr = append(arr, vv)
		}
	}
	return arr
}

func containsAny[T comparable](arr []T, v ...T) bool {
	for _, vv := range v {
		if slices.Contains(arr, vv) {
			return true
		}
	}
	return false
}

func glob(patterns []string) iter.Seq[string] {
	return func(yield func(string) bool) {
		visited := make(map[string]struct{})
		for _, pattern := range patterns {
			matches, err := filepath.Glob(pattern)
			if err != nil {
				log.Panicf("failed to glob: %s, %v", pattern, err)
			}
			for _, match := range matches {
				match, err = filepath.Rel(".", match)
				if err != nil {
					log.Panicf("failed to get relative path: %s, %v", match, err)
				}
				if _, ok := visited[match]; ok {
					continue
				}
				visited[match] = struct{}{}

				if !yield(match) {
					return
				}
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
	maxLen := 0
	for _, row := range csvData {
		if len(row) > maxLen {
			maxLen = len(row)
		}
	}
	for i, row := range csvData {
		if len(row) < maxLen {
			csvData[i] = append(row, make([]string, maxLen-len(row))...)
		}
	}

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

func has(arr any, v any) bool {
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() != reflect.Slice {
		panic("arr is not slice")
	}
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			if has(arr, value.Index(i).Interface()) {
				return true
			}
		}
	} else {
		value = value.Convert(arrValue.Type().Elem())
		for i := 0; i < arrValue.Len(); i++ {
			if arrValue.Index(i).Interface() == value.Interface() {
				return true
			}
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

func reflectContainsAny(a, b any) bool {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	if av.Kind() != reflect.Slice || bv.Kind() != reflect.Slice {
		panic("a or b is not slice")
	}
	for i := 0; i < av.Len(); i++ {
		for j := 0; j < bv.Len(); j++ {
			if av.Index(i).Interface() == bv.Index(j).Interface() {
				return true
			}
		}
	}
	return false
}

var pluralizeClient = pluralize.NewClient()

func singular(s string) string {
	return pluralizeClient.Singular(s)
}

func plural(s string) string {
	return pluralizeClient.Plural(s)
}

//go:embed templates/*
var templateFS embed.FS
var templateFuncMap = initTemplateFuncs()

func initTemplateFuncs() template.FuncMap {
	m := sprig.FuncMap()
	m["singular"] = singular
	m["plural"] = plural
	m["pascal"] = pascal
	m["has"] = has
	return m
}
