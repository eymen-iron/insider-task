package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/eymen-iron/insider-task/esrc"
	"github.com/eymen-iron/insider-task/router"
	"github.com/eymen-iron/insider-task/utils"
)

// Item structunu tanımla
type Item struct {
	ID       string `json:"item_id"`
	Name     string `json:"name"`
	Locale   string `json:"locale"`
	Click    int    `json:"click"`
	Purchase int    `json:"purchase"`
}

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		fmt.Println("Elasticsearch connection error:", err)
		return
	}

	jsonData, err := ioutil.ReadFile("sample.json")
	if err != nil {
		fmt.Println("Sample JSON doesn't read:", err)
		return
	}

	var items []esrc.Item
	err = json.Unmarshal(jsonData, &items)
	if err != nil {
		fmt.Println("JSON data doesn't index ", err)
		return
	}

	for _, item := range items {
		req := esapi.IndexRequest{
			Index:      "sample_index",
			DocumentID: item.ID,
			Body:       bytes.NewReader(utils.EncodeJSON(item)),
			Refresh:    "true",
		}

		res, err := req.Do(context.Background(), es)
		if err != nil {
			fmt.Println("Data indexing error", err)
			return
		}
		defer res.Body.Close()
	}

	// HTTP sunucusunu başlat
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		router.GetAllItems(w, r, es)
	})
	fmt.Println("HTTP Sunucusu Başlatıldı...")
	if err := http.ListenAndServe(":3001", nil); err != nil {
		fmt.Println("HTTP Sunucusu başlatılamadı:", err)
	}
}
