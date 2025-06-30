package restapi

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
			name:           "invalid JSON",
			requestBody:    `{invalid}`,
			mockError:      assert.AnError,
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
			// Create mock
			mockDeleter := new(MockURLDeleter)
			if tt.requestBody != "" {
				mockDeleter.On("DeleteURLS", mock.Anything, []byte(tt.requestBody)).Return(tt.mockError)
			}

			// Replace the real function with our mock
			originalDelete := usecasesDeleteURLS
			usecasesDeleteURLS = mockDeleter.DeleteURLS
			defer func() {
				usecasesDeleteURLS = originalDelete
			}()

			// Create request
			var body io.Reader
			if tt.requestBody != "" {
				body = bytes.NewBufferString(tt.requestBody)
			} else {
				body = nil
			}

			req, err := http.NewRequest("DELETE", "/api/user/urls", body)
			require.NoError(t, err)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler := HandlerAPIUserUrlsDel("http://example.com")
			handler(rr, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Verify mock expectations
			mockDeleter.AssertExpectations(t)
		})
	}
}
