package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockService struct{}

func (m mockService) GetOriginalURL(shortPart string) string      { return "mockOriginalURL" }
func (m mockService) GetShortURLPart(orig string) (string, error) { return "mockShortURLPart", nil }

func TestRootHandler(t *testing.T) {
	type have struct {
		contentType string
		body        string
	}
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		have have
		want want
	}{
		{
			name: "RootHandler test 1",
			have: have{
				contentType: "text/plain",
				body:        "http://ya.ru/",
			},
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name: "RootHandler test 2",
			have: have{
				contentType: "text/plain; charset=utf-8",
				body:        "http://ya.ru/",
			},
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name: "RootHandler test 3",
			have: have{
				contentType: "text/html; charset=utf-8",
				body:        "http://ya.ru/",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name: "RootHandler test 4",
			have: have{
				contentType: "",
				body:        "http://ya.ru/",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name: "RootHandler test 5",
			have: have{
				contentType: "text/plain; charset=utf-8",
				body:        "",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
	}

	handlers := NewHandlers("mockBaseURL", mockService{})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBody := strings.NewReader(test.have.body)
			r := httptest.NewRequest(http.MethodPost, "/", reqBody)
			r.Header.Set("Content-Type", test.have.contentType)
			w := httptest.NewRecorder()
			handlers.RootHandler(w, r)
			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), test.want.contentType)
			res.Body.Close()
		})
	}
}

func TestShortenHandler(t *testing.T) {
	type have struct {
		contentType string
		body        string
	}
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		have have
		want want
	}{
		{
			name: "ShortenHandler test 1",
			have: have{
				contentType: "application/json",
				body:        `{"url":"http://ya.ru"}`,
			},
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
		{
			name: "ShortenHandler test 2",
			have: have{
				contentType: "text/plain",
				body:        `{"url":"http://ya.ru"}`,
			},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "ShortenHandler test 3",
			have: have{
				contentType: "application/json",
				body:        `{}`,
			},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "ShortenHandler test 4",
			have: have{
				contentType: "application/json",
				body:        `{"url":""}`,
			},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "ShortenHandler test 5",
			have: have{
				contentType: "application/json",
				body:        `{"urls":"http://ya.ru"}`,
			},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	handlers := NewHandlers("mockBaseURL", mockService{})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBody := strings.NewReader(test.have.body)
			r := httptest.NewRequest(http.MethodPost, "/api/shorten", reqBody)
			r.Header.Set("Content-Type", test.have.contentType)
			w := httptest.NewRecorder()
			handlers.ShortenHandler(w, r)
			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), test.want.contentType)
			res.Body.Close()
		})
	}
}
