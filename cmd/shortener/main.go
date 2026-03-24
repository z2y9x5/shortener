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

	cnf := config.GetConfig()
	cnf.ApplyCLIArgs()
	cnf.ApplyEnvArgs()

	db := repository.NewMemoryRepository()
	shortener := service.NewShortener(db)
	handlers := handler.NewHandlers(cnf.BaseURL, shortener)

	mux := chi.NewRouter()
	mux.Post("/", logger.RequestLogger(handlers.RootHandler))
	mux.Get("/{id}", logger.RequestLogger(handlers.RootWithShortHandler))
	mux.Post("/api/shorten", logger.RequestLogger(handlers.ShortenHandler))

	if err := http.ListenAndServe(cnf.ServerAddr, mux); err != nil {
		log.Panic("error in ListenAndServe", zap.Error(err))
	}
}
