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
	Repo := handlers.GetRepo(&app)

	http.HandleFunc("/", Repo.Home)
	http.HandleFunc("/about", Repo.About)

	fmt.Printf("Serving on port %s\n", portNumber)
	err = http.ListenAndServe(portNumber, nil)

	if err != nil {
		fmt.Println("Failed to start the server")
		return
	}

}
