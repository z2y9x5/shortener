package repository

import (
	"errors"
)

// Интерфейс хранилища данных.
type Repository interface {
	Get(short string) string
	Put(short string, orig string) error
}

// Хранилище данных в памяти.
type memoryRepository struct {
	// Карта соответствия короткой части URL и оригинального URL.
	urlMap map[string]string
}

// Получить оригинальный URL по короткой части.
func (m *memoryRepository) Get(short string) string {
	return m.urlMap[short]
}

// Добавить короткую часть и оригинальный URL.
func (m *memoryRepository) Put(short string, orig string) error {
	if _, ok := m.urlMap[short]; ok {
		return errors.New("ключ уже существует")
	}
	m.urlMap[short] = orig
	return nil
}

// Конструктор хранилища данных в памяти.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		urlMap: make(map[string]string),
	}
}
