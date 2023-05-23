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
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/generals-quarters", handlers.Repo.GeneralsQuarters)
	mux.Get("/majors-suite", handlers.Repo.MajorsSuite)

	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/book-room", handlers.Repo.BookRoom)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	mux.Get("/login", handlers.Repo.ShowLogin)
	mux.Post("/login", handlers.Repo.PostShowLogin)
	mux.Get("/logout", handlers.Repo.Logout)

	// Any route that starts with /admin will be handled here
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)

		mux.Get("/reservations/{src}/{reservation_id}", handlers.Repo.AdminShowReservation)      // src can be "all" or "new"
		mux.Post("/reservations/{src}/{reservation_id}", handlers.Repo.PostAdminShowReservation) // src can be "all" or "new"
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
