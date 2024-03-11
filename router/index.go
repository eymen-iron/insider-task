package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
)

func GetAllItems(w http.ResponseWriter, r *http.Request, es *elasticsearch.Client) {
	// Query parameters
	var params struct {
		Size int    `json:"size"`
		Page int    `json:"page"`
		Sort string `json:"sort"`
	}

	params.Size, _ = strconv.Atoi(r.URL.Query().Get("size"))
	if params.Size <= 0 {
		params.Size = 1000 // Default size 1000
	}

	params.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if params.Page <= 0 {
		params.Page = 1 // Default page 1
	}
	params.Sort = r.URL.Query().Get("sort")

	// Elasticsearch query
	var query struct {
		Size  int `json:"size"`
		From  int `json:"from"`
		Query struct {
			MatchAll map[string]interface{} `json:"match_all"`
		} `json:"query"`
		Sort string `json:"sort"`
	}
	query.Size = params.Size
	query.From = (params.Page - 1) * params.Size
	query.Query.MatchAll = map[string]interface{}{}
	switch params.Sort {
	case "click":
		query.Sort = "click"
	case "purchase":
		query.Sort = "purchase"
	default:
		query.Sort = "_doc"
	}

	// Elasticsearch request
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		http.Error(w, "Request error", http.StatusInternalServerError)
		return
	}
	res, err := es.Search(
		es.Search.WithContext(r.Context()),
		es.Search.WithIndex("sample_index"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		http.Error(w, "Elasticsearch req error ", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Elasticsearch yanıtını parse et
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		http.Error(w, "Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}

	totalHits := int(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})

	documents := make([]map[string]interface{}, len(hits))
	for i, hit := range hits {
		documents[i] = hit.(map[string]interface{})["_source"].(map[string]interface{})
	}

	response := struct {
		Data      []map[string]interface{} `json:"data"`
		TotalHits int                      `json:"total_items"`
		Page      int                      `json:"page"`
		Size      int                      `json:"size"`
	}{
		Data:      documents,
		TotalHits: totalHits,
		Page:      params.Page,
		Size:      params.Size,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
