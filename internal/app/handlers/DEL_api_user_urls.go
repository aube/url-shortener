package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/aube/url-shortener/internal/app/ctxkeys"
)

const numWorkers int = 10

type StorageDelete interface {
	Delete(c context.Context, l []string) error
}

// HandlerAPIUserUrlsDel deletes multiple URLs for a user
// @Summary Delete user URLs
// @Description Deletes multiple shortened URLs belonging to the authenticated user
// @Tags URLs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body []string true "Array of URL hashes to delete"
// @Success 202 {string} string "Deletion request accepted"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user/urls [delete]
func HandlerAPIUserUrlsDel(store StorageDelete, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		var data []interface{}
		json.Unmarshal([]byte(body), &data)

		userID := r.Context().Value(ctxkeys.UserIDKey).(string)

		go asyncMagic(data, store, userID)

		w.WriteHeader(http.StatusAccepted)
	}
}

func asyncMagic(inputData []interface{}, store StorageDelete, userID string) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := generator(doneCh, inputData)
	channels := fanOut(doneCh, inputCh)
	resultCh := fanIn(doneCh, channels...)

	var hashes []string
	for hash := range resultCh {
		hashes = append(hashes, hash)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, ctxkeys.UserIDKey, userID)
	defer cancel()

	store.Delete(ctx, hashes)
}

func generator(doneCh chan struct{}, input []interface{}) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			var str string
			if s, ok := data.(string); ok {
				str = s
			}
			select {
			case <-doneCh:
				return
			case inputCh <- str:
			}
		}
	}()

	return inputCh
}

func add(doneCh chan struct{}, inputCh chan string) chan string {
	addRes := make(chan string)

	go func() {
		defer close(addRes)

		for data := range inputCh {
			select {
			case <-doneCh:
				return
			case addRes <- data:
			}
		}
	}()
	return addRes
}

func fanOut(doneCh chan struct{}, inputCh chan string) []chan string {
	channels := make([]chan string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		addResultCh := add(doneCh, inputCh)
		channels[i] = addResultCh
	}

	return channels
}

func fanIn(doneCh chan struct{}, resultChs ...chan string) chan string {
	finalCh := make(chan string)

	var wg sync.WaitGroup

	for _, ch := range resultChs {
		chClosure := ch

		// инкрементируем счётчик горутин, которые нужно подождать
		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-doneCh:
					return
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
