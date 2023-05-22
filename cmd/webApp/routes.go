package main

import (
	"net/http"

	"github.com/ASHWIN776/learning-Go/internal/config"
	"github.com/ASHWIN776/learning-Go/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(WriteToConsole)
	mux.Use(SessionLoad)
	mux.Use(Nosurf)

	// Routes
	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))
	mux.Get("/generals-quarters", http.HandlerFunc(handlers.Repo.GeneralsQuarters))
	mux.Get("/majors-suite", http.HandlerFunc(handlers.Repo.MajorsSuite))

	mux.Get("/make-reservation", http.HandlerFunc(handlers.Repo.MakeReservation))
	mux.Post("/make-reservation", http.HandlerFunc(handlers.Repo.PostReservation))
	mux.Get("/reservation-summary", http.HandlerFunc(handlers.Repo.ReservationSummary))

	mux.Get("/search-availability", http.HandlerFunc(handlers.Repo.SearchAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(handlers.Repo.AvailabilityJSON))
	mux.Get("/book-room", http.HandlerFunc(handlers.Repo.BookRoom))
	mux.Post("/search-availability", http.HandlerFunc(handlers.Repo.PostAvailability))
	mux.Get("/choose-room/{id}", http.HandlerFunc(handlers.Repo.ChooseRoom))

	mux.Get("/login", http.HandlerFunc(handlers.Repo.ShowLogin))
	mux.Post("/login", http.HandlerFunc(handlers.Repo.PostShowLogin))
	mux.Get("/logout", http.HandlerFunc(handlers.Repo.Logout))

	// Any route that starts with /admin will be handled here
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/dashboard", http.HandlerFunc(handlers.Repo.AdminDashboard))
	})

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
