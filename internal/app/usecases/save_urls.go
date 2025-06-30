package usecases

import (
	"context"
	"encoding/json"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StorageSetMultiple defines the interface for storing multiple URLs at once.
type StorageSetMultiple interface {
	// SetMultiple stores a map of hashes to URLs in storage.
	// Returns an error if the operation fails.
	SetMultiple(context.Context, map[string]string) error
}

// SaveURLS stores multiple URLs in batch and returns their shortened versions.
// The input is expected to be a JSON array of correlation IDs and original URLs.
// Returns a JSON array of correlation IDs with shortened URLs or an error if the operation fails.
func SaveURLS(ctx context.Context, body []byte, baseURL string) ([]byte, error) {
	store := store.NewStore()
	log := logger.WithContext(ctx)

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

	err := store.SetMultiple(ctx, items)

	if err != nil {
		return nil, err
	}

	log.Info("SaveMultipleURLs", "items length", len(items))

	return JSON2Batch(outputBatch)

}

// inputBatchJSONItem represents a single input item in the batch save operation.
type inputBatchJSONItem struct {
	ID  string `json:"correlation_id"` // The correlation ID for tracking
	URL string `json:"original_url"`   // The original URL to shorten
}

// outputBatchJSONItem represents a single output item in the batch save operation.
type outputBatchJSONItem struct {
	ID    string `json:"correlation_id"` // The correlation ID from input
	SHORT string `json:"short_url"`      // The generated shortened URL
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

// JSON2Batch json.Marshal
func JSON2Batch(outputJSON []outputBatchJSONItem) ([]byte, error) {
	log := logger.Get()
	jsonBytes, err := json.Marshal(outputJSON)

	if err != nil {
		log.Error("JSON2Batch", "err", err)
		return nil, err
	}

	return jsonBytes, nil
}
