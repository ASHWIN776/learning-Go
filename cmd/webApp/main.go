package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ASHWIN776/learning-Go/internal/config"
	"github.com/ASHWIN776/learning-Go/internal/handlers"
	"github.com/ASHWIN776/learning-Go/internal/models"
	"github.com/ASHWIN776/learning-Go/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}

	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Printf("Serving on port %s\n", portNumber)
	err = srv.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start the server")
		return
	}

}

func run() error {
	// What can I put in the session
	gob.Register(models.Reservation{})

	app.InProduction = false

	tc, err := render.BuildTemplateCache()

	if err != nil {
		return errors.New("cannot build template cache")
	}

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	app.TemplateCache = tc
	app.UseCache = false
	render.GetConfig(&app)

	// As Home and About are Methods of an instance of type Repository - So, we need the instance from render.go
	handlers.AddRepo(&app)

	return nil
}
