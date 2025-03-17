package handlers

import (
	"context"
	"errors"
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

type MockMemoryStore struct {
	s map[string]string
}

/* func (m *MockMemoryStore) Get(s string) (string, bool) {
	return s, true
} */

func (m *MockMemoryStore) Set(c context.Context, k string, v string) error {
	if v == "conflict" {
		return errors.New("")
	}
	return nil
}

/*
	 func (m *MockMemoryStore) List() map[string]string {
		return nil
	}

	func (m *MockMemoryStore) Ping() error {
		return nil
	}
*/

func TestHandlerAPI(t *testing.T) {
	baseURL := "http://localhost:8080"
	fakeAddress := "http://test.test/test"
	hash := hasher.CalcHash([]byte(fakeAddress))
	conflictHash := hasher.CalcHash([]byte("conflict"))

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

		{
			name: "conflict short URL",
			want: want{
				statusCode:   409,
				responseBody: baseURL + "/" + conflictHash,
			},
			postBody: "conflict",
		},

		{
			name: "error on empty body",
			want: want{
				statusCode:   400,
				responseBody: "Request body is empty\n",
			},
			postBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.postBody))
			w := httptest.NewRecorder()
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			h := HandlerRoot(ctx, MemoryStore, baseURL)
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
