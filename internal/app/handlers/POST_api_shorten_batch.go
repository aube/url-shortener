package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/logger"
)

type StorageSetMultiple interface {
	SetMultiple(context.Context, map[string]string) error
}

func HandlerShortenBatch(store StorageSetMultiple, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.WithContext(ctx)

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

		err = store.SetMultiple(r.Context(), items)

		if err != nil {
			log.Error("SetMultiple", "err", err)
			http.Error(w, "Failed to write URLs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(JSON2Batch(outputBatch))
	}
}

type inputBatchJSONItem struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

func batch2JSON(body []byte) []inputBatchJSONItem {
	log := logger.Get()

	inputJSON := []inputBatchJSONItem{}
	err := json.Unmarshal(body, &inputJSON)

	if err != nil {
		log.Error("batch2JSON", "err", err)
	}

	return inputJSON
}

type outputBatchJSONItem struct {
	ID    string `json:"correlation_id"`
	SHORT string `json:"short_url"`
}

func JSON2Batch(outputJSON []outputBatchJSONItem) []byte {
	log := logger.Get()
	jsonBytes, err := json.Marshal(outputJSON)

	if err != nil {
		log.Error("JSON2Batch", "err", err)
	}

	return jsonBytes
}
