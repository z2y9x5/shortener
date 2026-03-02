package config

// Значения по умолчанию.
const (
	DefaultProtocol = "http"
	DefaultHost     = "localhost"
	DefaultPort     = "8080"
)

// Конфигурация.
type Config struct {
	App App
}

// Конфигурация приложения.
type App struct {
	Host  string
	Port  string
	Proto string
}

// Конструктор конфигурации.
func NewConfig() *Config {
	return &Config{
		App: App{
			Host:  DefaultHost,
			Port:  DefaultPort,
			Proto: DefaultProtocol,
		},
	}
}
