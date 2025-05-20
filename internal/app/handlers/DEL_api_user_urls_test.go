package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockApi "github.com/aube/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlerAPIUserUrlsDel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mockApi.MockStorageDelete)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful deletion request",
			requestBody: `["abc123", "def456"]`,
			mockSetup: func(m *mockApi.MockStorageDelete) {
				m.EXPECT().
					Delete(gomock.Any(), []string{"abc123", "def456"}).
					Return(nil)
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:        "storage error",
			requestBody: `["abc123"]`,
			mockSetup: func(m *mockApi.MockStorageDelete) {
				m.EXPECT().
					Delete(gomock.Any(), []string{"abc123"}).
					Return(errors.New("storage error"))
			},
			expectedStatus: http.StatusAccepted, // Still accepts since operation is async
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mockApi.NewMockStorageDelete(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockStorage)
			}

			handler := HandlerAPIUserUrlsDel(mockStorage, "http://localhost")

			var req *http.Request
			if tt.requestBody != "" {
				req = httptest.NewRequest(http.MethodDelete, "/", strings.NewReader(tt.requestBody))
			} else {
				req = httptest.NewRequest(http.MethodDelete, "/", nil)
			}

			w := httptest.NewRecorder()
			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != "" {
				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}

func TestHandlerAPIUserUrlsDel_ContentType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockApi.NewMockStorageDelete(ctrl)
	mockStorage.EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(nil)

	handler := HandlerAPIUserUrlsDel(mockStorage, "http://localhost")
	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader(`["abc123"]`))
	w := httptest.NewRecorder()

	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusAccepted, res.StatusCode)
	assert.Equal(t, "", res.Header.Get("Content-Type")) // No content type set
}

func TestHandlerAPIUserUrlsDel_ContextPropagation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockApi.NewMockStorageDelete(ctrl)
	mockStorage.EXPECT().
		Delete(gomock.Any(), []string{"abc123"}).
		DoAndReturn(func(ctx context.Context, _ []string) error {
			if ctx == nil {
				t.Error("Expected non-nil context")
			}
			return nil
		})

	handler := HandlerAPIUserUrlsDel(mockStorage, "http://localhost")
	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader(`["abc123"]`))
	w := httptest.NewRecorder()

	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusAccepted, res.StatusCode)
}
