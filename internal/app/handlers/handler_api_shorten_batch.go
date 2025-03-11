package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/logger"
)

type StorageSetMultiple interface {
	SetMultiple(map[string]string) error
}

func HandlerShortenBatch(store StorageSetMultiple, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil || r.ContentLength == 0 {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Read the entire body content
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		inputJSON := batch2JSON(body)
		outputBatch := []outputBatchJSONItem{}
		items := make(map[string]string)

		for _, v := range inputJSON {
			hash := hasher.CalcHash([]byte(v.URL))
			outputBatch = append(outputBatch, outputBatchJSONItem{
				ID:    v.ID,
				SHORT: baseURL + "/" + hash,
			})
			items[hash] = v.URL
		}
		fmt.Println("items", items)
		fmt.Println("outputBatch", outputBatch)

		err = store.SetMultiple(items)
		fmt.Println("err", err)

		if err != nil {
			http.Error(w, "Failed to write URLs", http.StatusInternalServerError)
			return
		}

		// hash := hasher.CalcHash(originalURL)
		// store.Set(hash, string(originalURL))
		// shortURL := baseURL + "/" + hash

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `%s`, JSON2batch(outputBatch))

		// logger.Println("URL:", shortURL, http.StatusCreated)
	}
}

type inputBatchJSONItem struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}
type Target struct {
	Id    int     `json:"id"`
	Price float64 `json:"price"`
}

func batch2JSON(body []byte) []inputBatchJSONItem {

	inputJSON := []inputBatchJSONItem{}
	err := json.Unmarshal(body, &inputJSON)

	if err != nil {
		logger.Infoln(err)
	}

	return inputJSON
}

type outputBatchJSONItem struct {
	ID    string `json:"correlation_id"`
	SHORT string `json:"short_url"`
}

func JSON2batch(outputJSON []outputBatchJSONItem) []byte {
	jsonBytes, err := json.Marshal(outputJSON)

	if err != nil {
		logger.Infoln(err)
	}

	return jsonBytes
}

// type inputJSON struct {
// 	[]inputJSONItem
// }

// type outputJSON struct {
// 	[]outputJSONItem
// }

/* type JSONItem struct {
	Hash string `json:"short_url"`
	URL  string `json:"original_url"`
}

func urlsJSON(memData map[string]string, baseURL string) []byte {
	var jsonData []JSONItem

	for k, v := range memData {
		item := JSONItem{Hash: baseURL + "/" + k, URL: v}
		jsonData = append(jsonData, item)
	}
	jsonBytes, err := json.Marshal(jsonData)

	if err != nil {
		logger.Infoln(err)
	}

	return jsonBytes
}
*/
