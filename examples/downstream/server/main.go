package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/clienttable"
	"server/servertable"
)

type tableData struct {
	Data     []byte
	Checksum string
}

func mustJson(data interface{}) []byte {
	encoded, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return encoded
}

func main() {
	serverTables, err := servertable.LoadTablesFromFile("./tabledata")
	if err != nil {
		log.Panicf("failed to load server tables: %v", err)
	}

	log.Println("server table:")
	for _, r := range serverTables.SampleData.Rows {
		log.Printf("%+v\n", r)
	}

	clientTables, err := clienttable.LoadTablesFromFile("./tabledata")
	if err != nil {
		log.Panicf("failed to load client tables: %v", err)
	}

	fmt.Println("client table:")
	for _, r := range clientTables.SampleData.Rows {
		fmt.Printf("%+v\n", r)
	}

	tableCacheMap := make(map[string]tableData)
	for _, t := range clientTables.GetTables() {
		data := mustJson(t.GetRows())
		checksum := md5.Sum(data)
		tableCacheMap[t.SheetName()] = tableData{
			Data:     data,
			Checksum: hex.EncodeToString(checksum[:]),
		}
		log.Println("table cache:", t.SheetName(), tableCacheMap[t.SheetName()].Checksum)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /table/{sheetName}", func(res http.ResponseWriter, req *http.Request) {
		sheetName := req.PathValue("sheetName")
		checksum := req.URL.Query().Get("checksum")
		log.Printf("table req: %s, %s", sheetName, checksum)

		cache, ok := tableCacheMap[sheetName]
		if !ok {
			http.Error(res, "table not found", http.StatusNotFound)
			return
		}

		resBody := map[string]interface{}{
			"checksum": cache.Checksum,
		}
		if checksum != cache.Checksum {
			log.Printf("sending data: %s", sheetName)
			resBody["data"] = cache.Data
		}
		json.NewEncoder(res).Encode(resBody)
		res.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8081", mux)
}
