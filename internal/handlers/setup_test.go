package handlers

import (
	"encoding/gob"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ASHWIN776/learning-Go/internal/config"
	"github.com/ASHWIN776/learning-Go/internal/models"
	"github.com/ASHWIN776/learning-Go/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "../../templates"

func GetRoutes() http.Handler {
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	// Info log (will be a place to send all the info logs - currently to stdout)
	app.InfoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)

	// Error log (will be a place to send all the Error logs - currently to stdout)
	app.ErrorLog = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	tc, err := BuildTestTemplateCache()

	if err != nil {
		log.Println("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true
	render.GetConfig(&app)

	// As Home and About are Methods of an instance of type Repository - So, we need the instance from render.go
	AddRepo(&app)

	mux := chi.NewRouter()

	mux.Use(WriteToConsole)
	mux.Use(SessionLoad)
	// mux.Use(Nosurf)

	// Routes
	mux.Get("/", http.HandlerFunc(Repo.Home))
	mux.Get("/about", http.HandlerFunc(Repo.About))
	mux.Get("/contact", http.HandlerFunc(Repo.Contact))
	mux.Get("/generals-quarters", http.HandlerFunc(Repo.GeneralsQuarters))
	mux.Get("/majors-suite", http.HandlerFunc(Repo.MajorsSuite))

	mux.Get("/make-reservation", http.HandlerFunc(Repo.MakeReservation))
	mux.Post("/make-reservation", http.HandlerFunc(Repo.PostReservation))
	mux.Get("/reservation-summary", http.HandlerFunc(Repo.ReservationSummary))

	mux.Get("/search-availability", http.HandlerFunc(Repo.SearchAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(Repo.AvailabilityJSON))
	mux.Post("/search-availability", http.HandlerFunc(Repo.PostAvailability))

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

// middlewares

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Nosurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func BuildTestTemplateCache() (map[string]*template.Template, error) {

	tc := map[string]*template.Template{}

	// Find out all the pages ending with page.gohtml in the templates directory
	pages, pageErr := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))

	if pageErr != nil {
		log.Println("cannot find file in the specified path")
		return nil, errors.New("cannot build cache")
	}

	// Populate the cache
	for _, page := range pages {
		pageName := filepath.Base(page)
		var parseErr error
		tc[pageName], parseErr = template.ParseFiles(fmt.Sprintf("%s/%s", pathToTemplates, pageName), fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates))

		if parseErr != nil {
			log.Println("Parse Error: failed to parse ", pageName)
		}
	}

	// Return the cache
	return tc, nil
}
