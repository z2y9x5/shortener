package config

import (
	"flag"
)

// Значения по умолчанию.
const (
	defaultServerAddr = "localhost:8080"
	defaultBaseURL    = "http://" + defaultServerAddr
)

// Интерфейс конфигурации.
type Config interface {
	GetAppConfig() app
	ApplyCLIArgs()
}

// Конфигурация приложения.
type app struct {
	ServerAddr string
	BaseURL    string
}

// Конфигурация.
type config struct {
	App app
}

// Вернуть конфигурацию приложения.
func (c *config) GetAppConfig() app {
	return c.App
}

// Применить значения из аргументов командной строки.
func (c *config) ApplyCLIArgs() {
	flag.StringVar(&c.App.ServerAddr, "a", defaultServerAddr, "Адрес сервера в формате хост:порт. Пример: "+defaultServerAddr)
	flag.StringVar(&c.App.BaseURL, "b", defaultBaseURL, "Базовый URL для ссылок. Пример: "+defaultBaseURL)
	flag.Parse()
	deleteLastSlash(&c.App.BaseURL)
}

// Конструктор конфигурации.
func NewConfig() Config {
	return &config{
		App: app{
			ServerAddr: defaultServerAddr,
			BaseURL:    defaultBaseURL,
		},
	}
}

// Удалить слеш в конце строки.
func deleteLastSlash(url *string) {
	runes := []rune(*url)
	len := len(runes)
	if len > 0 && runes[len-1] == '/' {
		*url = string(runes[:len-1])
	}
}
