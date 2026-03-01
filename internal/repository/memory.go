package repository

import (
	"errors"

	"github.com/z2y9x5/shortener/internal/model"
)

// Хранилище данных в памяти.
type MemoryRepository struct {
	db model.URLMap
}

// Получить оригинальный URL по короткой части.
func (m *MemoryRepository) Get(short string) string {
	return m.db[short]
}

// Добавить короткую часть и оригинальный URL.
func (m *MemoryRepository) Put(short string, orig string) error {
	if _, ok := m.db[short]; ok {
		return errors.New("ключ уже существует")
	}
	m.db[short] = orig
	return nil
}

// Конструктор хранилища данных в памяти.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		db: model.URLMap{},
	}
}
