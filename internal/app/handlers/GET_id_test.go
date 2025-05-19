package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	mockApi "github.com/aube/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlerID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testID := "abc123"
	testURL := "http://example.com"

	tests := []struct {
		name           string
		id             string
		setupMock      func(*mockApi.MockStorageGet)
		expectedStatus int
		expectedHeader string
		expectedBody   string
	}{
		{
			name: "successful redirect",
			id:   testID,
			setupMock: func(m *mockApi.MockStorageGet) {
				m.EXPECT().
					Get(gomock.Any(), testID).
					Return(testURL, true)
			},
			expectedStatus: http.StatusTemporaryRedirect,
			expectedHeader: testURL,
		},
		{
			name: "URL not found",
			id:   testID,
			setupMock: func(m *mockApi.MockStorageGet) {
				m.EXPECT().
					Get(gomock.Any(), testID).
					Return("", false)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL not found\n",
		},
		{
			name: "URL deleted",
			id:   testID,
			setupMock: func(m *mockApi.MockStorageGet) {
				m.EXPECT().
					Get(gomock.Any(), testID).
					Return("", true)
			},
			expectedStatus: http.StatusGone,
			expectedBody:   "URL deleted\n",
		},
		{
			name:           "empty ID",
			id:             "",
			setupMock:      func(m *mockApi.MockStorageGet) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "ID must be specified\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mockApi.NewMockStorageGet(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStorage)
			}

			handler := HandlerID(mockStorage)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			if tt.id != "" {
				req.SetPathValue("id", tt.id)
			}
			w := httptest.NewRecorder()

			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedHeader != "" {
				assert.Equal(t, tt.expectedHeader, res.Header.Get("Location"))
			}

			if tt.expectedBody != "" {
				body := make([]byte, len(tt.expectedBody))
				_, err := res.Body.Read(body)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}

func TestHandlerID_ContextPropagation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testID := "abc123"
	testURL := "http://example.com"

	t.Run("passes context to storage", func(t *testing.T) {
		mockStorage := mockApi.NewMockStorageGet(ctrl)
		mockStorage.EXPECT().
			Get(gomock.Any(), testID).
			DoAndReturn(func(ctx context.Context, _ string) (string, bool) {
				// Verify context is passed through
				if ctx == nil {
					t.Error("Expected non-nil context")
				}
				return testURL, true
			})

		handler := HandlerID(mockStorage)
		req := httptest.NewRequest(http.MethodGet, "/"+testID, nil)
		req.SetPathValue("id", testID)
		w := httptest.NewRecorder()

		handler(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(t, testURL, res.Header.Get("Location"))
	})
}
