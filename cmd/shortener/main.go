package main

import (
	"log"
	"net/http"

	"github.com/z2y9x5/shortener/internal/config"
	"github.com/z2y9x5/shortener/internal/handler"

	"github.com/go-chi/chi"
)

func main() {
	cnf := config.NewConfig()

	// Инкремент 3. Маршрутизатор сторонней библиотеки chi.
	mux := chi.NewRouter()
	mux.Post("/", handler.RootHandler)
	mux.Get("/{id}", handler.RootWithShortHandler)

	// Инкремент 2. Маршрутизатор стандартной библиотеки net/http.
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handler.DefaultHandler)
	// mux.HandleFunc("POST /{$}", handler.RootHandler)
	// mux.HandleFunc("GET /{id}", handler.RootWithShortHandler)

	handler.ServerProtocol = cnf.App.Proto
	listen := cnf.App.Host + ":" + cnf.App.Port
	log.Fatal(http.ListenAndServe(listen, mux))
}
