// Code generated by "nestcsv"; DO NOT EDIT.

package {{ .PackageName }}

type TableBase interface {
    SheetName() string
    GetRows() interface{}
    Load(data []byte) error
    LoadFromString(jsonString string) error
    LoadFromFile(basePath string) error
}

