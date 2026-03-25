package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter реализует интерфейс [http.ResponseWriter] и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// Header - обертка для [http.ResponseWriter.Header].
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write - обертка для [http.ResponseWriter.Write].
func (c *compressWriter) Write(p []byte) (int, error) {
	if isCompressibleType(c.w.Header().Get("Content-Type")) {
		c.w.Header().Set("Content-Encoding", "gzip")
		return c.zw.Write(p)
	}
	c.w.Header().Del("Content-Encoding")
	return c.w.Write(p)
}

// WriteHeader - обертка для [http.ResponseWriter.WriteHeader].
func (c *compressWriter) WriteHeader(statusCode int) {
	if isCompressibleType(c.w.Header().Get("Content-Type")) {
		c.w.Header().Set("Content-Encoding", "gzip")
	} else {
		c.w.Header().Del("Content-Encoding")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает [gzip.Writer] и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	if isCompressibleType(c.w.Header().Get("Content-Type")) {
		return c.zw.Close()
	}
	return nil
}

// newCompressWriter - конструктор для [compressWriter].
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// разархивировать получаемые от клиента данные.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// Read - обертка для [io.ReadCloser.Read].
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close - обертка для [io.ReadCloser.Close].
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// newCompressReader - конструктор для [compressReader].
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// isCompressibleType проверяет "Content-Type".
// Сжимаемые типы: "application/json", "text/html".
func isCompressibleType(contentType string) bool {
	if strings.HasPrefix(contentType, "application/json") ||
		strings.HasPrefix(contentType, "text/html") {
		return true
	}
	return false
}

// GzipMiddleware поддерживает сжатые запросы с заголовком "Content-Encoding: gzip"
// и сжатые ответы при "Accept-Encoding: gzip".
func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origWriter := w

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			compressReader, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "error reading compressed data", http.StatusBadRequest)
				return
			}
			r.Body = compressReader
			defer compressReader.Close()
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			compressWriter := newCompressWriter(w)
			origWriter = compressWriter
			defer compressWriter.Close()
		}

		h.ServeHTTP(origWriter, r)
	})
}
