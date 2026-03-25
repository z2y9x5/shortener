package repository

import (
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/z2y9x5/shortener/internal/model"
)

// fileRepository хранит данные в файле.
type fileRepository struct {
	file    *os.File
	records *memoryRepository
}

// Close закрывает файл, предварительно записав данные в него.
// uuid генерируются порядковыми номерами и могут не совпадать с прошлой записью.
func (f *fileRepository) Close() error {
	defer f.file.Close()
	if err := f.file.Truncate(0); err != nil {
		return err
	}
	if _, err := f.file.Seek(0, 0); err != nil {
		return err
	}
	records := make([]model.FileRecord, 0)
	uuid := 1
	for k, v := range f.records.urlMap {
		records = append(records, model.FileRecord{
			Uuid:        strconv.Itoa(uuid),
			ShortUrl:    k,
			OriginalUrl: v,
		})
		uuid++
	}
	enc := json.NewEncoder(f.file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(records); err != nil {
		return err
	}
	return nil
}

// Get принимает короткую часть URL и возвращает оригинальный URL.
func (f *fileRepository) Get(short string) string {
	return f.records.Get(short)
}

// Put сохраняет короткую чать URL и оригинальный URL.
func (f *fileRepository) Put(short string, orig string) error {
	return f.records.Put(short, orig)
}

// NewFileRepository конструктор для [fileRepository].
// Если файл не существует, то он будет создан.
// Если файл существует, то данные будут считаны из него (при наличии).
// Атрибут uuid игнорируется.
func NewFileRepository(fileName string) (*fileRepository, error) {
	repo := &fileRepository{
		records: NewMemoryRepository(),
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return nil, err
	}
	repo.file = file
	info, err := repo.file.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() > 0 {
		data, err := io.ReadAll(repo.file)
		if err != nil {
			return nil, err
		}
		var fileRecords []model.FileRecord
		if err := json.Unmarshal(data, &fileRecords); err != nil {
			return nil, err
		}
		for _, r := range fileRecords {
			repo.records.Put(r.ShortUrl, r.OriginalUrl)
		}
	}
	return repo, nil
}
