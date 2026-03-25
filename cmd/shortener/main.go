package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	db, err := repository.NewFileRepository(cnf.URLFile)
	if err != nil {
		log.Panic("database error", zap.Error(err))
	}
	defer db.Close()

	shortener := service.NewShortener(db)
	handlers := handler.NewHandlers(cnf.BaseURL, shortener)

	mux := chi.NewRouter()
	mux.Use(logger.LoggerMiddleware)
	mux.Use(handler.GzipMiddleware)
	mux.Post("/", handlers.RootHandler)
	mux.Get("/{id}", handlers.RootWithShortHandler)
	mux.Post("/api/shorten", handlers.ShortenHandler)

	srv := &http.Server{
		Addr:    cnf.ServerAddr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic("error in ListenAndServe", zap.Error(err))
		}
	}()

	<-stop

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Panic("error in Shutdown", zap.Error(err))
	}
}
