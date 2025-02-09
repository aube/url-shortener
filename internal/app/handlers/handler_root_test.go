package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aube/url-shortener/internal/app/hashes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerRoot(t *testing.T) {
	baseUrl := "http://localhost:8080"
	fakeAddress := "http://test.test/test"

	hash := hashes.CalcHash([]byte(fakeAddress))
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
				shortAddress: baseUrl + "/" + hash,
			},
			postBody: fakeAddress,
		},

		{
			name: "empty body",
			want: want{
				statusCode:   400,
				shortAddress: "",
			},
			postBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.postBody))
			w := httptest.NewRecorder()
			h := func(w http.ResponseWriter, r *http.Request) {
				HandlerRoot(w, r, baseUrl)
			}
			h(w, r)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			shortAddressResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.shortAddress, string(shortAddressResult))
		})
	}
}
