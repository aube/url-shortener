package restapi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBatchSaver struct {
	mock.Mock
}

func (m *MockBatchSaver) SaveMultipleURLs(ctx context.Context, body []byte, baseURL string) ([]byte, error) {
	args := m.Called(ctx, body, baseURL)
	return args.Get(0).([]byte), args.Error(1)
}

func TestHandlerShortenBatch(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockResponse   []byte
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful batch",
			requestBody:    `[{"correlation_id":"1","original_url":"https://example.com"}]`,
			mockResponse:   []byte(`[{"correlation_id":"1","short_url":"http://test/abc123"}]`),
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "batch error",
			requestBody:    `[{"correlation_id":"1","original_url":"https://example.com"}]`,
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "empty body",
			requestBody:    "",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSaver := new(MockBatchSaver)
			if tt.requestBody != "" {
				mockSaver.On("SaveMultipleURLs", mock.Anything, []byte(tt.requestBody), "http://test").Return(tt.mockResponse, tt.mockError)
			}

			// Replace the real usecase with our mock
			originalSave := usecasesSaveURLS
			usecasesSaveURLS = mockSaver.SaveMultipleURLs
			defer func() { usecasesSaveURLS = originalSave }()

			req := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBufferString(tt.requestBody))
			rr := httptest.NewRecorder()

			handler := HandlerShortenBatch("http://test")
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.requestBody != "" {
				mockSaver.AssertExpectations(t)
			}
		})
	}
}
