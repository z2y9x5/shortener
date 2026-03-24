package logger

import (
	"bytes"
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
	body   *bytes.Buffer
}

// loggingResponseWriter подменяет исходный [http.ResponseWriter] в [RequestLogger].
type loggingResponseWriter struct {
	w        http.ResponseWriter
	respData *responseData
}

// Header - обертка для [http.ResponseWriter.Header].
func (l *loggingResponseWriter) Header() http.Header {
	return l.w.Header()
}

// Write - обертка для [http.ResponseWriter.Write].
func (l *loggingResponseWriter) Write(p []byte) (int, error) {
	l.respData.body.Write(p)
	count, err := l.w.Write(p)
	l.respData.size = count
	return count, err
}

// WriteHeader - обертка для [http.ResponseWriter.WriteHeader].
func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.respData.status = statusCode
	l.w.WriteHeader(statusCode)
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

		responseData := &responseData{
			body: &bytes.Buffer{},
		}
		lw := loggingResponseWriter{
			w:        w,
			respData: responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.Info("incoming request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
			zap.String("body", responseData.body.String()),
			zap.Duration("duration", duration),
		)
	})
}
