package usecases

import (
	"context"
	"encoding/json"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StorageSetMultiple interface
type StorageSetMultiple interface {
	SetMultiple(context.Context, map[string]string) error
}

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
