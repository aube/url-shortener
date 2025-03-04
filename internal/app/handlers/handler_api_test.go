package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockMemoryStore struct{}

func (m *MockMemoryStore) Get(s string) (string, bool) {
	return s, true
}

func (m *MockMemoryStore) Set(k string, v string) error {
	return nil
}

func TestHandlerAPI(t *testing.T) {
	baseURL := "http://localhost:8080"
	fakeAddress := "http://test.test/test"
	MemoryStore := &MockMemoryStore{}

	hash := hasher.CalcHash([]byte(fakeAddress))
	type want struct {
		statusCode   int
		shortAddress string
	}
	tests := []struct {
		name     string
		postBody string
		id       string
		want     want
	}{
		{
			name: "fakeAddress body",
			want: want{
				statusCode:   201,
				shortAddress: baseURL + "/" + hash,
			},
			postBody: fakeAddress,
		},

		{
			name: "empty body",
			want: want{
				statusCode:   400,
				shortAddress: "Request body is empty\n",
			},
			postBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.postBody))
			w := httptest.NewRecorder()
			h := HandlerRoot(MemoryStore, baseURL)
			h(w, r)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			shortAddressResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			err = result.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			require.NoError(t, err)

			assert.Equal(t, tt.want.shortAddress, string(shortAddressResult))
		})
	}
}
