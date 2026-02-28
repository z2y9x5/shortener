package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/z2y9x5/shortener/internal/service"
)

// Протокол по умолчанию.
var ServerProtocol = "http"

// Эндпоинт по умолчанию.
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

// Эндпоинт с методом POST и путём /.
// Сервер принимает в теле запроса строку URL как text/plain.
// Возвращает ответ с кодом 201 и сокращённым URL как text/plain.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if !(strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain")) {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	shortURLPart, err := service.GetShortURLPart(string(body))
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusBadRequest)
		return
	}
	resp := ServerProtocol + "://" + r.Host + "/" + shortURLPart
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(resp))
}

// Эндпоинт с методом GET и путём /{id}, где id — идентификатор сокращённого URL.
// Cервер возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func RootWithShortHandler(w http.ResponseWriter, r *http.Request) {
	shortPart := r.PathValue("id")
	originalURL := service.GetOriginalURL(shortPart)
	if originalURL == "" {
		http.Error(w, "The requested URL was not found", http.StatusBadRequest)
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
