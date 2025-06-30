package restapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aube/url-shortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockURLGetter is a mock implementation of URLGetter interface
type MockURLGetter struct {
	mock.Mock
}

func (m *MockURLGetter) GetURL(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func TestHandlerID(t *testing.T) {

	tests := []struct {
		name           string
		id             string
		mockURL        string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful redirect",
			id:             "abc123",
			mockURL:        "https://example.com",
			mockError:      nil,
			expectedStatus: http.StatusTemporaryRedirect,
			expectedBody:   "",
		},
		{
			name:           "empty ID",
			id:             "",
			mockURL:        "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "ID must be specified\n",
		},
		{
			name:           "URL not found",
			id:             "notfound",
			mockURL:        "",
			mockError:      assert.AnError,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL not found\n",
		},
		{
			name:           "URL deleted",
			id:             "deleted",
			mockURL:        "",
			mockError:      nil,
			expectedStatus: http.StatusGone,
			expectedBody:   "URL deleted\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock URL getter
			mockGetter := new(MockURLGetter)

			// Set up mock expectations only when ID is not empty
			if tt.id != "" {
				mockGetter.On("GetURL", mock.Anything, tt.id).Return(tt.mockURL, tt.mockError)
			}

			// Create handler with our mock getter
			handler := NewHandlerID(mockGetter)

			// Create request with path value
			req, err := http.NewRequest("GET", "/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set path value if ID is not empty
			if tt.id != "" {
				req.SetPathValue("id", tt.id)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check response body
			assert.Equal(t, tt.expectedBody, rr.Body.String())

			// For successful redirect, check Location header
			if tt.expectedStatus == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.mockURL, rr.Header().Get("Location"))
			}

			// Assert mock expectations
			if tt.id != "" {
				mockGetter.AssertExpectations(t)
			}
		})
	}
}

// NewHandlerID creates a new HandlerID with dependency injection
func NewHandlerID(getter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.WithContext(ctx)

		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "ID must be specified", http.StatusBadRequest)
			return
		}

		log.Debug("HandlerID", "id", id)

		url, err := getter.GetURL(ctx, id)

		if err != nil {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		if url == "" {
			http.Error(w, "URL deleted", http.StatusGone)
			return
		}

		log.Debug("HandlerID", "url", url)

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// URLGetter defines the interface for getting URLs
type URLGetter interface {
	GetURL(ctx context.Context, id string) (string, error)
}
