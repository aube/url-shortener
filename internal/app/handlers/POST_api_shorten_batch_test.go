package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockApi "github.com/aube/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlerShortenBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseURL := "http://localhost:8080"

	tests := []struct {
		name             string
		requestBody      string
		setupMock        func(*mockApi.MockStorageSetMultiple)
		expectedStatus   int
		expectedResponse string
		expectErrorLog   bool
	}{
		{
			name: "storage error during batch creation",
			requestBody: `[
				{"correlation_id": "1", "original_url": "http://example.com"}
			]`,
			setupMock: func(m *mockApi.MockStorageSetMultiple) {
				m.EXPECT().SetMultiple(gomock.Any(), map[string]string{
					"89dce6a446": "http://example.com",
				}).Return(errors.New("storage error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "Failed to write URLs\n",
			expectErrorLog:   true,
		},
		{
			name:             "empty request body",
			requestBody:      "",
			setupMock:        func(m *mockApi.MockStorageSetMultiple) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "Request body is empty\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mockApi.NewMockStorageSetMultiple(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStorage)
			}

			handler := HandlerShortenBatch(mockStorage, baseURL)

			var req *http.Request
			if tt.requestBody != "" {
				req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			} else {
				req = httptest.NewRequest(http.MethodPost, "/", nil)
			}

			w := httptest.NewRecorder()
			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedResponse != "" {
				body := make([]byte, len(tt.expectedResponse))
				_, err := res.Body.Read(body)
				assert.NoError(t, err)

				if strings.HasPrefix(tt.expectedResponse, "[") {
					fmt.Println(string(body))
					assert.JSONEq(t, strings.TrimSpace(tt.expectedResponse), strings.TrimSpace(string(body)))
				} else {
					assert.Equal(t, tt.expectedResponse, string(body))
				}
			}

			// Verify Content-Type header for successful responses
			if res.StatusCode == http.StatusCreated {
				assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestBatch2JSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []inputBatchJSONItem
	}{
		{
			name: "valid JSON array",
			input: `[
				{"correlation_id": "1", "original_url": "http://example.com"},
				{"correlation_id": "2", "original_url": "http://test.org"}
			]`,
			expected: []inputBatchJSONItem{
				{ID: "1", URL: "http://example.com"},
				{ID: "2", URL: "http://test.org"},
			},
		},
		{
			name:     "empty array",
			input:    `[]`,
			expected: []inputBatchJSONItem{},
		},
		{
			name:     "invalid JSON",
			input:    `not json`,
			expected: []inputBatchJSONItem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := batch2JSON([]byte(tt.input))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJSON2Batch(t *testing.T) {
	tests := []struct {
		name     string
		input    []outputBatchJSONItem
		expected string
	}{
		{
			name: "valid output items",
			input: []outputBatchJSONItem{
				{ID: "1", SHORT: "http://localhost/hash1"},
				{ID: "2", SHORT: "http://localhost/hash2"},
			},
			expected: `[
				{"correlation_id": "1", "short_url": "http://localhost/hash1"},
				{"correlation_id": "2", "short_url": "http://localhost/hash2"}
			]`,
		},
		{
			name:     "empty array",
			input:    []outputBatchJSONItem{},
			expected: `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JSON2Batch(tt.input)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}
