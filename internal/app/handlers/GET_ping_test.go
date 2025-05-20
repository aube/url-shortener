package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockApi "github.com/aube/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlerPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMock      func(*mockApi.MockStoragePing)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful ping",
			setupMock: func(m *mockApi.MockStoragePing) {
				m.EXPECT().
					Ping(gomock.Any()).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "pong",
		},
		{
			name: "failed ping",
			setupMock: func(m *mockApi.MockStoragePing) {
				m.EXPECT().
					Ping(gomock.Any()).
					Return(errors.New("connection failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL not found\n",
		},
		{
			name: "context canceled",
			setupMock: func(m *mockApi.MockStoragePing) {
				m.EXPECT().
					Ping(gomock.Any()).
					Return(context.Canceled)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mockApi.NewMockStoragePing(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStorage)
			}

			handler := HandlerPing(mockStorage)

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != "" {
				body := make([]byte, len(tt.expectedBody))
				_, err := res.Body.Read(body)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}

func TestHandlerPing_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	errorCases := []struct {
		name        string
		err         error
		expectError bool
	}{
		{"database error", errors.New("db timeout"), true},
		{"context deadline", context.DeadlineExceeded, true},
		{"nil error", nil, false},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStorage := mockApi.NewMockStoragePing(ctrl)
			mockStorage.EXPECT().
				Ping(gomock.Any()).
				Return(tc.err)

			handler := HandlerPing(mockStorage)
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if tc.expectError {
				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			} else {
				assert.Equal(t, http.StatusOK, res.StatusCode)
			}
		})
	}
}
