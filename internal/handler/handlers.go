package handler

import (
	"io"
	"net/http"
	"strings"

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
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	shortURLPart, err := h.shortener.GetShortURLPart(string(body))
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusBadRequest)
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
		http.Error(w, "The requested URL was not found", http.StatusBadRequest)
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// NewHandlers конструктор для [handlers].
func NewHandlers(url string, service Shortener) *handlers {
	return &handlers{
		shortener: service,
		baseURL:   url,
	}
}
