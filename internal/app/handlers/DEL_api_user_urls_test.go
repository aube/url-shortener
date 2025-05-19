package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aube/url-shortener/internal/app/ctxkeys"
	mockApi "github.com/aube/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlerAPIUserUrlsDel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserID := "user123"
	testHashes := []string{"abc123", "def456"}
	jsonInput, _ := json.Marshal(testHashes)

	tests := []struct {
		name           string
		requestBody    string
		userContext    interface{}
		setupMock      func(*mockApi.MockStorageDelete)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "successful deletion request",
			requestBody:    string(jsonInput),
			userContext:    testUserID,
			setupMock:      func(m *mockApi.MockStorageDelete) {},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "invalid JSON format",
			requestBody:    "not json",
			userContext:    testUserID,
			setupMock:      func(m *mockApi.MockStorageDelete) {},
			expectedStatus: http.StatusAccepted, // Still accepts since processing is async
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mockApi.NewMockStorageDelete(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStorage)
			}

			handler := HandlerAPIUserUrlsDel(mockStorage, "http://localhost")

			var req *http.Request
			if tt.requestBody != "" {
				req = httptest.NewRequest(http.MethodDelete, "/", strings.NewReader(tt.requestBody))
			} else {
				req = httptest.NewRequest(http.MethodDelete, "/", nil)
			}

			ctx := context.WithValue(req.Context(), ctxkeys.UserIDKey, tt.userContext)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedError {
				body := make([]byte, 100)
				_, err := res.Body.Read(body)
				assert.NoError(t, err)
				assert.Contains(t, string(body), "Failed")
			}
		})
	}
}

func TestAsyncMagic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserID := "user123"
	testHashes := []interface{}{"abc123", "def456"}

	t.Run("processes valid input correctly", func(t *testing.T) {
		mockStorage := mockApi.NewMockStorageDelete(ctrl)
		mockStorage.EXPECT().
			Delete(gomock.Any(), []string{"abc123", "def456"}).
			DoAndReturn(func(ctx context.Context, hashes []string) error {
				// Verify context has user ID
				if ctx.Value(ctxkeys.UserIDKey) != testUserID {
					t.Error("User ID not set in context")
				}
				return nil
			})

		asyncMagic(testHashes, mockStorage, testUserID)
	})
}

func TestFanOutFanIn(t *testing.T) {
	t.Run("processes all inputs through workers", func(t *testing.T) {
		doneCh := make(chan struct{})
		defer close(doneCh)

		inputCh := make(chan string)
		go func() {
			defer close(inputCh)
			inputCh <- "a"
			inputCh <- "b"
			inputCh <- "c"
		}()

		channels := fanOut(doneCh, inputCh)
		resultCh := fanIn(doneCh, channels...)

		var results []string
		for val := range resultCh {
			results = append(results, val)
		}

		assert.ElementsMatch(t, []string{"a", "b", "c"}, results)
	})

	t.Run("handles early termination", func(t *testing.T) {
		doneCh := make(chan struct{})
		inputCh := make(chan string)

		go func() {
			defer close(inputCh)
			inputCh <- "a"
			close(doneCh)
			inputCh <- "b" // Should be ignored
		}()

		channels := fanOut(doneCh, inputCh)
		resultCh := fanIn(doneCh, channels...)

		var results []string
		for val := range resultCh {
			results = append(results, val)
		}

		assert.LessOrEqual(t, len(results), 1) // Only "a" might get through
	})
}
