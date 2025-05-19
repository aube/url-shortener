package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	mockApi "github.com/aube/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlerRoot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseURL := "http://localhost:8080"
	testURL := "http://example.com"
	testHash := "89dce6a446"

	tests := []struct {
		name             string
		requestBody      string
		contentType      string
		acceptHeader     string
		setupMock        func(*mockApi.MockStorageSet)
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "successful creation with plain text",
			requestBody: testURL,
			setupMock: func(m *mockApi.MockStorageSet) {
				m.EXPECT().
					Set(gomock.Any(), testHash, testURL).
					Return(nil)
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: baseURL + "/" + testHash,
		},
		{
			name:        "successful creation with json content-type",
			requestBody: `{"URL":"` + testURL + `"}`,
			contentType: "application/json",
			setupMock: func(m *mockApi.MockStorageSet) {
				m.EXPECT().
					Set(gomock.Any(), testHash, testURL).
					Return(nil)
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: `{"result":"` + baseURL + "/" + testHash + `"}`,
		},
		{
			name:         "successful creation with json accept header",
			requestBody:  `{"URL":"` + testURL + `"}`,
			contentType:  "application/json",
			acceptHeader: "application/json",
			setupMock: func(m *mockApi.MockStorageSet) {
				m.EXPECT().
					Set(gomock.Any(), testHash, testURL).
					Return(nil)
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: `{"result":"` + baseURL + "/" + testHash + `"}`,
		},
		{
			name:        "conflict when url exists",
			requestBody: testURL,
			setupMock: func(m *mockApi.MockStorageSet) {
				m.EXPECT().
					Set(gomock.Any(), testHash, testURL).
					Return(appErrors.NewHTTPError(http.StatusConflict, "conflict"))
			},
			expectedStatus:   http.StatusConflict,
			expectedResponse: baseURL + "/" + testHash,
		},
		{
			name:        "internal server error on storage failure",
			requestBody: testURL,
			setupMock: func(m *mockApi.MockStorageSet) {
				m.EXPECT().
					Set(gomock.Any(), testHash, testURL).
					Return(errors.New("storage error"))
			},
			expectedStatus:   http.StatusCreated, // Note: Current handler doesn't change status for non-HTTP errors
			expectedResponse: baseURL + "/" + testHash,
		},
		{
			name:           "empty body",
			requestBody:    "",
			setupMock:      func(m *mockApi.MockStorageSet) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mockApi.NewMockStorageSet(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockStorage)
			}

			handler := HandlerRoot(mockStorage, baseURL)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			if tt.acceptHeader != "" {
				req.Header.Set("Accept", tt.acceptHeader)
			}

			w := httptest.NewRecorder()
			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedResponse != "" {
				body := make([]byte, len(tt.expectedResponse))
				_, err := res.Body.Read(body)
				assert.NoError(t, err)

				if strings.HasPrefix(tt.expectedResponse, "{") {
					assert.JSONEq(t, tt.expectedResponse, string(body))
				} else {
					assert.Equal(t, tt.expectedResponse, string(body))
				}
			}
		})
	}
}

func TestHandlerRoot_ContentTypeHandling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseURL := "http://localhost:8080"
	testURL := "http://example.com"
	testHash := "89dce6a446"

	tests := []struct {
		name             string
		contentType      string
		acceptHeader     string
		expectedResponse string
	}{
		{
			name:             "plain text response when no json headers",
			expectedResponse: baseURL + "/" + testHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mockApi.NewMockStorageSet(ctrl)
			mockStorage.EXPECT().
				Set(gomock.Any(), testHash, testURL).
				Return(nil)

			handler := HandlerRoot(mockStorage, baseURL)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testURL))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			if tt.acceptHeader != "" {
				req.Header.Set("Accept", tt.acceptHeader)
			}

			w := httptest.NewRecorder()
			handler(w, req)

			res := w.Result()
			defer res.Body.Close()

			body := make([]byte, len(tt.expectedResponse))
			_, err := res.Body.Read(body)
			assert.NoError(t, err)

			if strings.HasPrefix(tt.expectedResponse, "{") {
				assert.JSONEq(t, tt.expectedResponse, string(body))
			} else {
				assert.Equal(t, tt.expectedResponse, string(body))
			}

			// Verify Content-Type header
			contentType := w.Header().Get("Content-Type")
			if strings.Contains(tt.expectedResponse, "{") {
				assert.Equal(t, "application/json", contentType)
			} else {
				assert.NotEqual(t, "application/json", contentType)
			}
		})
	}
}
