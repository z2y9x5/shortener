package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log хранит общий для приложения логер.
var Log *zap.Logger = zap.NewNop()

// responseData хранит информацию об ответе на запрос.
// Используется в [loggingResponseWriter].
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter подменяет исходный [http.ResponseWriter] в [RequestLogger].
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Initialize принимает уровень логирования в виде строки.
func Initialize(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	Log = logger
	return Log, nil
}

// LoggerMiddleware - middleware-логер входящих запросов.
// Логирует URI запроса, метод запроса, код статуса ответа,
// размер содержимого ответа, время выполнения.
func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Log.Sync()
		start := time.Now()

		responseData := &responseData{}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.Info("incoming request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
			zap.Duration("duration", duration),
		)
	})
}
