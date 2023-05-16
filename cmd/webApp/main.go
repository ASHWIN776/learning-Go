package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ASHWIN776/learning-Go/internal/config"
	"github.com/ASHWIN776/learning-Go/internal/driver"
	"github.com/ASHWIN776/learning-Go/internal/handlers"
	"github.com/ASHWIN776/learning-Go/internal/helpers"
	"github.com/ASHWIN776/learning-Go/internal/models"
	"github.com/ASHWIN776/learning-Go/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	// Defer the connection close
	defer db.SQL.Close()

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

func run() (*driver.DB, error) {
	// What can I put in the session
	gob.Register(models.Reservation{})

	app.InProduction = false

	tc, err := render.BuildTemplateCache()

	if err != nil {
		return nil, errors.New("cannot build template cache")
	}

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// Info log (will be a place to send all the info logs - currently to stdout)
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// Error log (will be a place to send all the Error logs - currently to stdout)
	errorLog := log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	app.Session = session
	app.TemplateCache = tc
	app.UseCache = false

	// Connecting to database
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=pass@3750")

	if err != nil {
		return nil, err
	}
	log.Println("Connected to Database")

	render.GetConfig(&app)
	helpers.GetConfig(&app)

	// As Home and About are Methods of an instance of type Repository - So, we need the instance from render.go
	handlers.AddRepo(&app, db)

	return db, nil
}
