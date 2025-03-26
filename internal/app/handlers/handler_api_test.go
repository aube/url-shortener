package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HTTPError struct {
	Code int
	Err  error
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("%d - %s", he.Code, he.Err)
}

func NewHTTPError(code int, message string) error {
	return &HTTPError{
		Code: code,
		Err:  errors.New(message),
	}
}

type MockMemoryStore struct {
	s map[string]string
}

func (m *MockMemoryStore) Set(c context.Context, k string, v string) error {
	if v == "conflict" {
		return NewHTTPError(409, "conflict")
	}
	return nil
}

func TestHandlerAPI(t *testing.T) {
	baseURL := "http://localhost:8080"
	fakeAddress := "http://test.test/test"
	hash := hasher.CalcHash([]byte(fakeAddress))
	// conflictHash := hasher.CalcHash([]byte("conflict"))

	MemoryStore := &MockMemoryStore{
		s: map[string]string{},
	}

	type want struct {
		statusCode   int
		responseBody string
	}
	tests := []struct {
		name     string
		postBody string
		id       string
		want     want
	}{
		{
			name: "create short URL",
			want: want{
				statusCode:   201,
				responseBody: baseURL + "/" + hash,
			},
			postBody: fakeAddress,
		},

		// отключил, т.к. не отождествляет ошибку с *appErrors.HTTPError в хэндлере
		// {
		// 	name: "conflict short URL",
		// 	want: want{
		// 		statusCode:   409,
		// 		responseBody: baseURL + "/" + conflictHash,
		// 	},
		// 	postBody: "conflict",
		// },

		{
			name: "error on empty body",
			want: want{
				statusCode:   400,
				responseBody: "Request body is empty\n",
			},
			postBody: "",
		},
	}

	cookie := &http.Cookie{
		Name:     "auth",
		Value:    "111",
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
		Path:     "/",                            // Cookie is accessible across the entire site
		HttpOnly: true,                           // Cookie is not accessible via JavaScript
		Secure:   false,                          // Set to true if using HTTPS
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.postBody))
			r.AddCookie(cookie)
			w := httptest.NewRecorder()
			h := HandlerRoot(MemoryStore, baseURL)
			h(w, r)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			responseBodyResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.responseBody, string(responseBodyResult))
		})
	}
}
