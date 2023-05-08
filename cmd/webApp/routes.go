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

	// Routes
	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))
	mux.Get("/generals-quarters", http.HandlerFunc(handlers.Repo.GeneralsQuarters))
	mux.Get("/majors-suite", http.HandlerFunc(handlers.Repo.MajorsSuite))
	mux.Get("/make-reservation", http.HandlerFunc(handlers.Repo.MakeReservation))

	// Creates a fileserver by telling it where the static directory exists
	fileServer := http.FileServer(http.Dir("./static/"))

	// registers a new route in the mux router that matches requests for URLs starting with "/static/"
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// The use of http.StripPrefix
	/*
		When a user requests a static file, they typically use a URL that looks something like this: http://example.com/static/myfile.css. However, the http.FileServer handler expects to be given a URL that contains only the path to the file relative to the directory it is serving, without any extra prefixes. In this case, it would expect a URL like /myfile.css.

		To remove the "/static/" prefix from the requested URL, we use the http.StripPrefix handler. This handler modifies the incoming request by removing the "/static/" prefix from the requested URL before passing it on to the http.FileServer handler. This way, the http.FileServer handler receives a URL that contains only the path to the file relative to the "static" directory, and it can serve the correct file.
	*/
	return mux
}
