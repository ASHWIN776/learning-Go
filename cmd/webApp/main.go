package main

import (
	"fmt"
	"log"
	"myApp/pkg/config"
	"myApp/pkg/handlers"
	"myApp/pkg/render"
	"net/http"
)

const portNumber = ":8000"

func main() {

	var app config.AppConfig
	tc, err := render.BuildTemplateCache()

	if err != nil {
		log.Fatal("cannot build Template Cache")
	}

	app.TemplateCache = tc
	render.GetConfig(&app)

	// As Home and About are Methods of an instance of type Repository - So, we need the instance from render.go
	handlers.AddRepo(&app)

	fmt.Printf("Serving on port %s\n", portNumber)

	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start the server")
		return
	}

}
