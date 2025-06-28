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

type MockURLDeleter struct {
	mock.Mock
}

func (m *MockURLDeleter) DeleteURLS(ctx context.Context, body []byte) error {
	args := m.Called(ctx, body)
	return args.Error(0)
}

func TestHandlerAPIUserUrlsDel(t *testing.T) {

	tests := []struct {
		name           string
		requestBody    string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			requestBody:    `["url1", "url2"]`,
			mockError:      nil,
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "empty body",
			requestBody:    "",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "deletion error",
			requestBody:    `["url1", "url2"]`,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDeleter := new(MockURLDeleter)
			if tt.requestBody != "" {
				mockDeleter.On("DeleteURLS", mock.Anything, []byte(tt.requestBody)).Return(tt.mockError)
			}

			// Replace the real usecase with our mock
			originalDelete := usecasesDeleteURLS
			usecasesDeleteURLS = mockDeleter.DeleteURLS
			defer func() { usecasesDeleteURLS = originalDelete }()

			req := httptest.NewRequest("DELETE", "/api/user/urls", bytes.NewBufferString(tt.requestBody))
			rr := httptest.NewRecorder()

			handler := HandlerAPIUserUrlsDel("http://test")
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.requestBody != "" {
				mockDeleter.AssertExpectations(t)
			}
		})
	}
}
