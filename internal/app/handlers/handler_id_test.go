package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func TestHandlerID(t *testing.T) {

	fakeAddress := "http://test.test/test"
	hash := hasher.CalcHash([]byte(fakeAddress))

	MemoryStore := store.NewMemoryStore()
	MemoryStore.Set(hash, fakeAddress)

	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		request string
		id      string
		want    want
	}{
		{
			name: hash,
			want: want{
				statusCode: 307,
				location:   fakeAddress,
			},
			request: "/",
			id:      hash,
		},

		{
			name: "long string",
			want: want{
				statusCode: 400,
				location:   "",
			},
			request: "/alongfakestringalongfakestringalongfakestringalongfakestring",
			id:      "alongfakestringalongfakestringalongfakestringalongfakestring",
		},
		{
			name: "empty string",
			want: want{
				statusCode: 400,
				location:   "",
			},
			request: "/",
			id:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			r.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandlerID)
			h(w, r)

			result := w.Result()

			err := result.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}
}
