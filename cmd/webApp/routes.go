package main

import (
	"net/http"

	"github.com/ASHWIN776/learning-Go/pkg/config"
	"github.com/ASHWIN776/learning-Go/pkg/handlers"

	"github.com/go-chi/chi/v5"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(WriteToConsole)
	mux.Use(SessionLoad)

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	return mux
}
