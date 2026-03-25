package repository

import (
	"errors"
	"sync"
)

// memoryRepository хранит данные в памяти.
type memoryRepository struct {
	sync.Mutex
	urlMap map[string]string // Ключ - короткая часть, значение - оригинальный URL.
}

// Get принимает короткую часть URL и возвращает оригинальный URL.
func (m *memoryRepository) Get(short string) string {
	m.Lock()
	res := m.urlMap[short]
	m.Unlock()
	return res
}

// Put сохраняет короткую чать URL и оригинальный URL.
func (m *memoryRepository) Put(short string, orig string) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.urlMap[short]; ok {
		return errors.New("ключ уже существует")
	}
	m.urlMap[short] = orig
	return nil
}

// NewMemoryRepository конструктор для [memoryRepository].
func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{
		urlMap: make(map[string]string),
	}
}
