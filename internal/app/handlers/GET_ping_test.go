package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPinger struct {
	mock.Mock
}

func (m *MockPinger) StorePing(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHandlerPing(t *testing.T) {
	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful ping",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "pong",
		},
		{
			name:           "ping error",
			mockError:      assert.AnError,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPinger := new(MockPinger)
			mockPinger.On("StorePing", mock.Anything).Return(tt.mockError)

			// Replace the real usecase with our mock
			originalPing := usecasesStorePing
			usecasesStorePing = mockPinger.StorePing
			defer func() { usecasesStorePing = originalPing }()

			req := httptest.NewRequest("GET", "/ping", nil)
			rr := httptest.NewRecorder()

			handler := HandlerPing()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			mockPinger.AssertExpectations(t)
		})
	}
}
