package config

import (
	"flag"
	"os"
)

const (
	defaultServerAddr = "localhost:8080"              // Прослушиваемый сервером адрес и порт.
	defaultBaseURL    = "http://" + defaultServerAddr // Базовый URL для короткой ссылки.
)

// cnf - общий для приложения экземпляр конфигурации.
// Доступен через вызов [GetConfig].
var cnf *Config

// Config хранит конфигурацию приложения.
type Config struct {
	ServerAddr string
	BaseURL    string
}

// ApplyCLIArgs меняет значения в [Config] на значения из флагов командной строки.
func (c *Config) ApplyCLIArgs() {
	flag.StringVar(&c.ServerAddr, "a", defaultServerAddr, "Адрес сервера в формате хост:порт. Пример: "+defaultServerAddr)
	flag.StringVar(&c.BaseURL, "b", defaultBaseURL, "Базовый URL для ссылок. Пример: "+defaultBaseURL)
	flag.Parse()
	c.deleteLastSlash(&c.BaseURL)
}

// ApplyEnvArgs меняет значения в [Config] на значения из переменных окружения.
func (c *Config) ApplyEnvArgs() {
	if val := os.Getenv("SERVER_ADDRESS"); val != "" {
		c.ServerAddr = val
	}
	if val := os.Getenv("BASE_URL"); val != "" {
		c.BaseURL = val
		c.deleteLastSlash(&c.BaseURL)
	}
}

// deleteLastSlash удаляет слеш в конце адреса [Config.BaseURL].
func (c *Config) deleteLastSlash(url *string) {
	runes := []rune(*url)
	len := len(runes)
	if len > 0 && runes[len-1] == '/' {
		*url = string(runes[:len-1])
	}
}

// GetConfig возвращает экземпляр конфигурации.
func GetConfig() *Config {
	if cnf == nil {
		return &Config{
			ServerAddr: defaultServerAddr,
			BaseURL:    defaultBaseURL,
		}
	}
	return cnf
}
