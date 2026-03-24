package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/z2y9x5/shortener/internal/model"

	"github.com/go-chi/chi"
)

// Shortener - интерфейс сервиса коротких URL.
type Shortener interface {
	GetOriginalURL(shortPart string) string
	GetShortURLPart(orig string) (string, error)
}

// handlers - эндпоинты сервера.
type handlers struct {
	shortener Shortener
	baseURL   string
}

// RootHandler - эндпоинт с методом POST и путём /.
// Сервер принимает в теле запроса строку URL как text/plain.
// Возвращает ответ с кодом 201 и сокращённым URL как text/plain.
func (h handlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	if !(strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain")) {
		http.Error(w, "content-type must be text/plain", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}
	shortURLPart, err := h.shortener.GetShortURLPart(string(body))
	if err != nil {
		http.Error(w, "failed to create short url", http.StatusBadRequest)
		return
	}
	resp := h.baseURL + "/" + shortURLPart
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(resp))
}

// RootWithShortHandler - эндпоинт с методом GET и путём /{id}, где id — идентификатор сокращённого URL.
// Cервер возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func (h handlers) RootWithShortHandler(w http.ResponseWriter, r *http.Request) {
	shortPart := chi.URLParam(r, "id")
	originalURL := h.shortener.GetOriginalURL(shortPart)
	if originalURL == "" {
		http.Error(w, "the requested url was not found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// ShortenHandler принимает в теле запроса JSON-объект [model.JSONRequest]
// и возвращает в ответ объект [model.JSONResponse].
func (h handlers) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if !(strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")) {
		http.Error(w, "content-type must be application/json", http.StatusBadRequest)
		return
	}
	var jsonReq model.JSONRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&jsonReq); err != nil {
		http.Error(w, "cannot decode request json body", http.StatusBadRequest)
		return
	}
	if jsonReq.URL == "" {
		http.Error(w, "url in json is empty", http.StatusBadRequest)
		return
	}
	shortURLPart, err := h.shortener.GetShortURLPart(jsonReq.URL)
	if err != nil {
		http.Error(w, "failed to create short url", http.StatusBadRequest)
		return
	}
	resp := model.JSONResponse{
		Result: h.baseURL + "/" + shortURLPart,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		http.Error(w, "error encoding response", http.StatusBadRequest)
		return
	}
}

// NewHandlers конструктор для [handlers].
func NewHandlers(url string, service Shortener) *handlers {
	return &handlers{
		shortener: service,
		baseURL:   url,
	}
}
