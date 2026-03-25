package service

import (
	"errors"
	"math/rand"
)

const (
	// Допустимые символы в короткой части ссылки.
	symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"
	// Длина короткой части ссылки.
	shortPartLength = 8
	// Максимальное число попыток сгенерировать короткую часть.
	shortPartMaxAttempts = 10
)

// Repository - интерфейс хранилища данных.
type Repository interface {
	Get(short string) string
	Put(short string, orig string) error
}

// shortener - сервис укорачивания url.
type shortener struct {
	repo Repository
}

// GetOriginalURL возвращает оригинальный URL по короткой части.
func (s *shortener) GetOriginalURL(shortPart string) string {
	return s.repo.Get(shortPart)
}

// GetShortURLPart генерирует короткую ссылку и помещает ее в хранилище.
// Количество попыток сгенерировать уникальную ссылку задается в [shortPartMaxAttempts].
// В случае неудачи возвращает пустую строку и ошибку.
func (s *shortener) GetShortURLPart(orig string) (string, error) {
	for range shortPartMaxAttempts {
		short := s.generateShortPart()
		if err := s.repo.Put(short, orig); err == nil {
			return short, nil
		}
	}
	return "", errors.New("попытки сгенерировать короткую часть исчерпаны")
}

// generateShortPart формирует короткую ссылку.
// Длина ссылки задается [shortPartyLength].
// Допустимые символы хранятся в [symbols].
func (s *shortener) generateShortPart() string {
	var short [shortPartLength]rune
	sym := []rune(symbols)
	for i := range shortPartLength {
		rnd := rand.Intn(len(sym))
		short[i] = sym[rnd]
	}
	return string(short[:])
}

// NewShortener конструктор для [shortener].
func NewShortener(r Repository) *shortener {
	return &shortener{
		repo: r,
	}
}
