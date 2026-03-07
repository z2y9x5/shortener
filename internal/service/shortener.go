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

// Интерфейс сервиса коротких url.
type Shortener interface {
	GetOriginalURL(shortPart string) string
	GetShortURLPart(orig string) (string, error)
}

// Сервис укорачивания url.
type shortener struct {
	// Хранилище URL.
	repo repository.Repository
}

// Получить оригинальный URL по короткой части.
func (s *shortener) GetOriginalURL(shortPart string) string {
	return s.repo.Get(shortPart)
}

// Получить короткую часть URL.
func (s *shortener) GetShortURLPart(orig string) (string, error) {
	for range shortPartMaxAttempts {
		short := s.generateShortPart()
		if err := s.repo.Put(short, orig); err == nil {
			return short, nil
		}
	}
	return "", errors.New("попытки сгенерировать короткую часть исчерпаны")
}

// Сгенерировать короткую часть.
func (s *shortener) generateShortPart() string {
	var short [shortPartLength]rune
	sym := []rune(symbols)
	for i := range shortPartLength {
		rnd := rand.Intn(len(sym))
		short[i] = sym[rnd]
	}
	return string(short[:])
}

// Конструктор сервиса коротких url.
func NewShortener(r repository.Repository) Shortener {
	return &shortener{
		repo: r,
	}
}
