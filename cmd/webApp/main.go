package main

import (
	"fmt"
	"log"
	"myApp/pkg/config"
	"myApp/pkg/handlers"
	"myApp/pkg/render"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8000"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false

	tc, err := render.BuildTemplateCache()

	if err != nil {
		log.Fatal("cannot build Template Cache")
	}

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	app.TemplateCache = tc
	render.GetConfig(&app)

	// As Home and About are Methods of an instance of type Repository - So, we need the instance from render.go
	handlers.AddRepo(&app)

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
