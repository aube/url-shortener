package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/logger"
)

type StorageSetMultiple interface {
	SetMultiple(context.Context, map[string]string) error
}

func HandlerShortenBatch(ctx context.Context, store StorageSetMultiple, baseURL string) http.HandlerFunc {
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

		err = store.SetMultiple(ctx, items)
		fmt.Println("err", err)

		if err != nil {
			http.Error(w, "Failed to write URLs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `%s`, JSON2batch(outputBatch))
	}
}

type inputBatchJSONItem struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
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
