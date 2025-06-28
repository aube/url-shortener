package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLSaverAPI struct {
	mock.Mock
}

func (m *MockURLSaverAPI) SaveURL(ctx context.Context, originalURL []byte, baseURL string) (string, error) {
	args := m.Called(ctx, originalURL, baseURL)
	return args.String(0), args.Error(1)
}

func TestHandlerAPI(t *testing.T) {

	tests := []struct {
		name           string
		requestBody    string
		mockShortURL   string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful creation",
			requestBody:    `{"url":"https://example.com"}`,
			mockShortURL:   "http://test/abc123",
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "conflict",
			requestBody:    `{"url":"https://example.com"}`,
			mockShortURL:   "http://test/abc123",
			mockError:      assert.AnError,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "empty body",
			requestBody:    "",
			mockShortURL:   "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSaver := new(MockURLSaverAPI)
			if tt.requestBody != "" {
				originalURL := "https://example.com" // Simplified for test
				mockSaver.On("SaveURL", mock.Anything, []byte(originalURL), "http://test").Return(tt.mockShortURL, tt.mockError)
			}

			// Replace the real usecase with our mock
			originalSave := usecasesSaveURL
			usecasesSaveURL = mockSaver.SaveURL
			defer func() { usecasesSaveURL = originalSave }()

			req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBufferString(tt.requestBody))
			rr := httptest.NewRecorder()

			handler := HandlerAPI("http://test")
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.requestBody != "" {
				mockSaver.AssertExpectations(t)
			}
		})
	}
}
