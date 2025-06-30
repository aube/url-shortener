package restapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockGetURLS mocks the usecases.GetURLS function
type MockGetURLS struct {
	mock.Mock
}

func (m *MockGetURLS) GetURLS(ctx context.Context, baseURL string) ([]byte, error) {
	args := m.Called(ctx, baseURL)
	return args.Get(0).([]byte), args.Error(1)
}

func TestHandlerAPIUserUrls(t *testing.T) {

	baseURL := "http://localhost:8080"

	tests := []struct {
		name           string
		mockReturnJSON []byte
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
		expectHeaders  map[string]string
	}{
		{
			name:           "Success with URLs",
			mockReturnJSON: []byte(`[{"short_url":"http://localhost:8080/abc123","original_url":"https://example.com"}]`),
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"short_url":"http://localhost:8080/abc123","original_url":"https://example.com"}]`,
			expectHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:           "Success with empty array",
			mockReturnJSON: []byte(""),
			mockReturnErr:  nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
			expectHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:           "Unauthorized error",
			mockReturnJSON: nil,
			mockReturnErr:  appErrors.NewHTTPError(http.StatusUnauthorized, "unauthorized"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
			expectHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:           "Internal server error",
			mockReturnJSON: nil,
			mockReturnErr:  errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to getJSON\n",
			expectHeaders: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
		},
		{
			name:           "Invalid JSON data",
			mockReturnJSON: []byte("{invalid}"),
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "{invalid}",
			expectHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockGetURLS := new(MockGetURLS)
			mockGetURLS.On("GetURLS", mock.Anything, baseURL).Return(tt.mockReturnJSON, tt.mockReturnErr)

			// Replace the real usecase with our mock
			originalUsecase := usecasesGetURLS
			usecasesGetURLS = mockGetURLS.GetURLS
			defer func() { usecasesGetURLS = originalUsecase }()

			// Create request and recorder
			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			rec := httptest.NewRecorder()

			// Call handler
			handler := HandlerAPIUserUrls(baseURL)
			handler(rec, req)

			// Verify results
			assert.Equal(t, tt.expectedStatus, rec.Code, "status code mismatch")

			if tt.expectedBody == "" {
				assert.Empty(t, rec.Body.String(), "expected empty body")
			} else {
				assert.Equal(t, tt.expectedBody, rec.Body.String(), "body content mismatch")
			}

			for key, value := range tt.expectHeaders {
				assert.Equal(t, value, rec.Header().Get(key), "header mismatch for "+key)
			}

			// Verify mock expectations
			mockGetURLS.AssertExpectations(t)
		})
	}
}

func TestGetJSON(t *testing.T) {
	tests := []struct {
		name     string
		memData  map[string]string
		baseURL  string
		expected string
		wantErr  bool
	}{
		{
			name: "Single URL",
			memData: map[string]string{
				"abc123": "https://example.com",
			},
			baseURL:  "http://localhost:8080",
			expected: `[{"short_url":"http://localhost:8080/abc123","original_url":"https://example.com"}]`,
			wantErr:  false,
		},
		{
			name: "Multiple URLs",
			memData: map[string]string{
				"abc123": "https://example.com",
				"def456": "https://another.com",
			},
			baseURL: "http://localhost:8080",
			expected: `[
				{"short_url":"http://localhost:8080/abc123","original_url":"https://example.com"},
				{"short_url":"http://localhost:8080/def456","original_url":"https://another.com"}
			]`,
			wantErr: false,
		},
		{
			name: "URL with special characters",
			memData: map[string]string{
				"gh789": "https://example.com/path?query=value",
			},
			baseURL:  "http://localhost:8080",
			expected: `[{"short_url":"http://localhost:8080/gh789","original_url":"https://example.com/path?query=value"}]`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getJSON(tt.memData, tt.baseURL)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(got))
		})
	}
}

func TestHandlerAPIUserUrls_EdgeCases(t *testing.T) {

	baseURL := "http://localhost:8080"

	t.Run("Nil context", func(t *testing.T) {
		// Setup mock
		mockGetURLS := new(MockGetURLS)
		mockGetURLS.On("GetURLS", mock.Anything, baseURL).
			Return([]byte("[]"), nil)

		// Replace the real usecase with our mock
		originalUsecase := usecasesGetURLS
		usecasesGetURLS = mockGetURLS.GetURLS
		defer func() { usecasesGetURLS = originalUsecase }()

		// Create request with nil context
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		req = req.WithContext(context.TODO())
		rec := httptest.NewRecorder()

		// Call handler
		handler := HandlerAPIUserUrls(baseURL)
		handler(rec, req)

		// Verify results
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockGetURLS.AssertExpectations(t)
	})
}
