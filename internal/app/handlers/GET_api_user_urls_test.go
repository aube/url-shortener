package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	mockApi "github.com/aube/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlerAPIUserUrls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseURL := "http://localhost:8080"
	testData := map[string]string{
		"abc123": "http://example.com",
		"def456": "http://test.org",
	}

	tests := []struct {
		name           string
		setupMock      func(*mockApi.MockStorageList)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful response with URLs",
			setupMock: func(m *mockApi.MockStorageList) {
				m.EXPECT().
					List(gomock.Any()).
					Return(testData, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `[
				{"short_url":"http://localhost:8080/abc123","original_url":"http://example.com"},
				{"short_url":"http://localhost:8080/def456","original_url":"http://test.org"}
			]`,
		},
		{
			name: "empty response",
			setupMock: func(m *mockApi.MockStorageList) {
				m.EXPECT().
					List(gomock.Any()).
					Return(map[string]string{}, nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "HTTP error from storage",
			setupMock: func(m *mockApi.MockStorageList) {
				m.EXPECT().
					List(gomock.Any()).
					Return(nil, appErrors.NewHTTPError(http.StatusInternalServerError, "storage error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "generic error from storage",
			setupMock: func(m *mockApi.MockStorageList) {
				m.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("generic error"))
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mockApi.NewMockStorageList(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStorage)
			}

			handler := HandlerAPIUserUrls(mockStorage, baseURL)

			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

			if tt.expectedBody != "" {
				var expected, actual []JSONItem
				err := json.Unmarshal([]byte(tt.expectedBody), &expected)
				assert.NoError(t, err)

				body := make([]byte, 1024)
				n, err := res.Body.Read(body)
				assert.NoError(t, err)

				err = json.Unmarshal(body[:n], &actual)
				assert.NoError(t, err)

				assert.ElementsMatch(t, expected, actual)
			}
		})
	}
}

func TestGetJSON(t *testing.T) {
	baseURL := "http://localhost:8080"
	testCases := []struct {
		name     string
		input    map[string]string
		expected []JSONItem
	}{
		{
			name: "single URL",
			input: map[string]string{
				"abc123": "http://example.com",
			},
			expected: []JSONItem{
				{Hash: baseURL + "/abc123", URL: "http://example.com"},
			},
		},
		{
			name: "multiple URLs",
			input: map[string]string{
				"abc123": "http://example.com",
				"def456": "http://test.org",
			},
			expected: []JSONItem{
				{Hash: baseURL + "/abc123", URL: "http://example.com"},
				{Hash: baseURL + "/def456", URL: "http://test.org"},
			},
		},
		{
			name:     "empty input",
			input:    map[string]string{},
			expected: []JSONItem{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getJSON(tc.input, baseURL)
			assert.NoError(t, err)

			var actual []JSONItem
			err = json.Unmarshal(result, &actual)
			assert.NoError(t, err)

			assert.ElementsMatch(t, tc.expected, actual)
		})
	}
}

func TestHandlerAPIUserUrls_EdgeCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("handles JSON marshal error", func(t *testing.T) {
		// This test would require mocking json.Marshal
		// Typically we'd use an interface for json operations to make this testable
		// Current implementation doesn't handle this error case explicitly
		t.Skip("JSON marshal error handling not implemented in current code")
	})

	t.Run("context propagation", func(t *testing.T) {
		mockStorage := mockApi.NewMockStorageList(ctrl)
		mockStorage.EXPECT().
			List(gomock.Any()).
			DoAndReturn(func(ctx context.Context) (map[string]string, error) {
				if ctx == nil {
					t.Error("Expected non-nil context")
				}
				return map[string]string{}, nil
			})

		handler := HandlerAPIUserUrls(mockStorage, "http://localhost")
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})
}
