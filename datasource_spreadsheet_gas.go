package nestcsv

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GASOption struct {
	URL                  string   `yaml:"url"`
	Password             string   `yaml:"password"`
	GoogleDriveFolderIDs []string `yaml:"google_drive_folder_ids"`
	SpreadsheetFileIDs   []string `yaml:"spreadsheet_file_ids"`
	DebugSaveDir         *string  `yaml:"debug_save_dir,omitempty"`

	// TODO : add google oauth2 authentication
}

func CollectSpreadsheetsThroughGAS(out chan<- *TableData, option *GASOption) error {
	zipData, err := callGASAndReadBase64(option)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("failed to read the zip: %w", err)
	}

	ch := make(chan *zip.File, 1000)
	go func() {
		for _, zipFile := range zipReader.File {
			if strings.HasPrefix(zipFile.Name, "#") {
				continue
			}
			ch <- zipFile
		}
		close(ch)
	}()

	var wg errgroup.Group
	for zipFile := range ch {
		wg.Go(func() error {
			file, err := zipFile.Open()
			if err != nil {
				return fmt.Errorf("failed to open the file: %s, %w", zipFile.Name, err)
			}
			defer file.Close()

			rows, err := csv.NewReader(file).ReadAll()
			if err != nil {
				return fmt.Errorf("failed to read the file: %s, %w", zipFile.Name, err)
			}

			tableData, err := ParseTableData(zipFile.Name, rows)
			if err != nil {
				if errors.Is(err, ErrSkipTable) {
					return nil
				}
				return err
			}
			if option.DebugSaveDir != nil {
				if err := tableData.SaveAsCSV(*option.DebugSaveDir); err != nil {
					return err
				}
			}
			out <- tableData
			return nil
		})
	}
	return wg.Wait()
}

func callGASAndReadBase64(option *GASOption) ([]byte, error) {
	uri, err := url.Parse(option.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	queryValues := url.Values{
		"password":  {option.Password},
		"folderIds": option.GoogleDriveFolderIDs,
		"fileIds":   option.SpreadsheetFileIDs,
	}
	uri.RawQuery = queryValues.Encode()

	res, err := http.Get(uri.String())
	if err != nil {
		return nil, fmt.Errorf("failed to download zip: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download zip: %s", res.Status)
	}
	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body : %w", err)
	}

	if len(resBytes) == 0 {
		return nil, fmt.Errorf("empty response body")
	}
	if resBytes[0] == '<' {
		return nil, fmt.Errorf("html response error: %s", string(resBytes))
	}

	maxDecodedLen := base64.StdEncoding.DecodedLen(len(resBytes))
	decoded := make([]byte, maxDecodedLen)
	n, err := base64.StdEncoding.Decode(decoded, resBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return decoded[:n], nil
}
