package main

import (
	"client/table"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const (
	serverUrl    = "http://localhost:8081"
	saveDataPath = "./savedata/"
)

func main() {
	var tables table.Tables
	for _, t := range tables.GetTables() {
		data, checksum := loadFile(t.SheetName())
		url := serverUrl + "/table/" + t.SheetName() + "?checksum=" + checksum
		log.Printf("req: %s\n", url)
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			panic("failed to get table: " + res.Status)
		}

		var resBody map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			panic(err)
		}

		// If the checksum is different, update the file.
		if checksum != resBody["checksum"].(string) {
			decoded, err := base64.StdEncoding.DecodeString(resBody["data"].(string))
			if err != nil {
				panic(err)
			}
			data = decoded
			if err := os.WriteFile(saveDataPath+t.SheetName()+".json", decoded, 0644); err != nil {
				panic(err)
			}
			log.Printf("update table from server: %s, %s\n", t.SheetName(), string(data))
		} else {
			log.Printf("table is up to date: %s\n", t.SheetName())
		}
		if err := t.Load(data); err != nil {
			panic(err)
		}
	}

	log.Println("client table:")
	for _, row := range tables.SampleData.Rows {
		log.Printf("%+v\n", row)
	}
}

func loadFile(sheetName string) ([]byte, string) {
	data, err := os.ReadFile(saveDataPath + sheetName + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ""
		}
		panic(err)
	}
	log.Println("load file:", sheetName, string(data))
	checksum := md5.Sum(data)
	return data, hex.EncodeToString(checksum[:])
}
