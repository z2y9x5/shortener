package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/z2y9x5/shortener/internal/config"
	"github.com/z2y9x5/shortener/internal/repository"
	"github.com/z2y9x5/shortener/internal/service"
)

func TestRootHandler(t *testing.T) {
	cnf := config.NewConfig()
	cnfApp := cnf.GetAppConfig()
	db := repository.NewMemoryRepository()
	shortener := service.NewShortener(db)
	handlers := NewHandlers(cnfApp.BaseURL, shortener)

	type want struct {
		requestContentType  string
		requestBody         string
		responseCode        int
		responseContentType string
		responseBody        string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "RootHandler test 1",
			want: want{
				requestContentType:  "text/plain; charset=utf-8",
				requestBody:         "https://ya.ru/",
				responseCode:        201,
				responseContentType: "text/plain",
				responseBody:        cnfApp.BaseURL,
			},
		},
		{
			name: "RootHandler test 2",
			want: want{
				requestContentType:  "text/html; charset=utf-8",
				requestBody:         "https://ya.ru/",
				responseCode:        400,
				responseContentType: "text/plain",
				responseBody:        "",
			},
		},
		{
			name: "RootHandler test 3",
			want: want{
				requestContentType:  "text/plain; charset=utf-8",
				requestBody:         "",
				responseCode:        400,
				responseContentType: "text/plain",
				responseBody:        "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBody := strings.NewReader(test.want.requestBody)
			r := httptest.NewRequest(http.MethodPost, "/", reqBody)
			r.Host = cnfApp.BaseURL
			r.Header.Set("Content-Type", test.want.requestContentType)
			w := httptest.NewRecorder()
			handlers.RootHandler(w, r)
			res := w.Result()

			assert.Equal(t, test.want.responseCode, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), test.want.responseContentType)
			resBody, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
			assert.Contains(t, string(resBody), test.want.responseBody)
		})
	}
}
