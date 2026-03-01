package service

import (
	"errors"
	"math/rand"

	"github.com/z2y9x5/shortener/internal/repository"
)

const (
	// Допустимые символы в короткой части ссылки.
	symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"
	// Длина короткой части ссылки.
	shortPartLength = 8
	// Максимальное число попыток сгенерировать короткую часть.
	shortPartMaxAttempts = 10
)

var (
	// Хранилище URL.
	repo URLRepository = repository.NewMemoryRepository()
)

// Интерфейс репозитория URL.
type URLRepository interface {
	Get(short string) string             // Получить оригинальный URL по короткой части.
	Put(short string, orig string) error // Добавить короткую часть и оригинальный URL.
}

// Сгенерировать короткую часть.
func generateShortPart() string {
	var short [shortPartLength]rune
	sym := []rune(symbols)
	for i := range shortPartLength {
		rnd := rand.Intn(len(sym))
		short[i] = sym[rnd]
	}
	return string(short[:])
}

// Получить короткую часть URL.
func GetShortURLPart(orig string) (string, error) {
	for range shortPartMaxAttempts {
		short := generateShortPart()
		if err := repo.Put(short, orig); err == nil {
			return short, nil
		}
	}
	return "", errors.New("попытки сгенерировать короткую часть исчерпаны")
}

// Получить оригинальный URL по короткой части.
func GetOriginalURL(shortPart string) string {
	return repo.Get(shortPart)
}
