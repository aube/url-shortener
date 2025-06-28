package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/app/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRootURLSaverRoot struct {
	mock.Mock
}

func (m *MockRootURLSaverRoot) SaveURL(ctx context.Context, originalURL []byte, baseURL string) (string, error) {
	args := m.Called(ctx, originalURL, baseURL)
	return args.String(0), args.Error(1)
}

func TestHandlerRoot(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		contentType    string
		acceptHeader   string
		mockShortURL   string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful text creation",
			requestBody:    "https://example.com",
			contentType:    "text/plain",
			acceptHeader:   "",
			mockShortURL:   "http://test/abc123",
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "successful json creation",
			requestBody:    `{"url":"https://example.com"}`,
			contentType:    "application/json",
			acceptHeader:   "",
			mockShortURL:   "http://test/abc123",
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "http error",
			requestBody:    "https://example.com",
			contentType:    "text/plain",
			acceptHeader:   "",
			mockShortURL:   "",
			mockError:      &appErrors.HTTPError{Code: http.StatusBadRequest},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			requestBody:    "",
			contentType:    "text/plain",
			acceptHeader:   "",
			mockShortURL:   "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSaver := new(MockRootURLSaverRoot)
			if tt.requestBody != "" {
				var originalURL string
				if tt.contentType == "application/json" {
					originalURL = "https://example.com" // Simplified for test
				} else {
					originalURL = tt.requestBody
				}
				mockSaver.On("SaveURL", mock.Anything, []byte(originalURL), "http://test").Return(tt.mockShortURL, tt.mockError)
			}

			// Replace the real usecase with our mock
			originalSave := usecases.SaveURL
			usecasesSaveURL = mockSaver.SaveURL
			defer func() { usecasesSaveURL = originalSave }()

			req := httptest.NewRequest("POST", "/", bytes.NewBufferString(tt.requestBody))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			if tt.acceptHeader != "" {
				req.Header.Set("Accept", tt.acceptHeader)
			}
			rr := httptest.NewRecorder()

			handler := HandlerRoot("http://test")
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.requestBody != "" {
				mockSaver.AssertExpectations(t)
			}
		})
	}
}
