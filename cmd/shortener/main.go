package main

import (
	"net/http"

	"github.com/z2y9x5/shortener/internal/config"
	"github.com/z2y9x5/shortener/internal/handler"
	"github.com/z2y9x5/shortener/internal/logger"
	"github.com/z2y9x5/shortener/internal/repository"
	"github.com/z2y9x5/shortener/internal/service"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.Initialize("info")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cnf := config.GetConfig()
	cnf.ApplyCLIArgs()
	cnf.ApplyEnvArgs()

	db := repository.NewMemoryRepository()
	shortener := service.NewShortener(db)
	handlers := handler.NewHandlers(cnf.BaseURL, shortener)

	mux := chi.NewRouter()
	mux.Use(logger.LoggerMiddleware)
	mux.Use(handler.GzipMiddleware)
	mux.Post("/", handlers.RootHandler)
	mux.Get("/{id}", handlers.RootWithShortHandler)
	mux.Post("/api/shorten", handlers.ShortenHandler)

	if err := http.ListenAndServe(cnf.ServerAddr, mux); err != nil {
		log.Panic("error in ListenAndServe", zap.Error(err))
	}
}
