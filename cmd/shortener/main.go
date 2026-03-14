package main

import (
	"log"
	"net/http"

	"github.com/z2y9x5/shortener/internal/config"
	"github.com/z2y9x5/shortener/internal/handler"
	"github.com/z2y9x5/shortener/internal/repository"
	"github.com/z2y9x5/shortener/internal/service"

	"github.com/go-chi/chi"
)

func main() {
	cnf := config.NewConfig()
	cnf.ApplyCLIArgs()
	cnf.ApplyEnvArgs()
	cnfApp := cnf.GetAppConfig()

	db := repository.NewMemoryRepository()
	shortener := service.NewShortener(db)
	handlers := handler.NewHandlers(cnfApp.BaseURL, shortener)

	mux := chi.NewRouter()
	mux.Post("/", handlers.RootHandler)
	mux.Get("/{id}", handlers.RootWithShortHandler)

	log.Fatal(http.ListenAndServe(cnfApp.ServerAddr, mux))
}
