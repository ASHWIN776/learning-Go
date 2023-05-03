package render

import (
	"errors"
	"html/template"
	"log"
	"myApp/pkg/config"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// Gets the pointer to the app config created in main.go, and assigns it to the pointer of the same type created here
func GetConfig(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, file string) {
	// Get the template cache from the App config
	tc := app.TemplateCache

	execErr := tc[file].Execute(w, nil)

	if execErr != nil {
		log.Println(execErr)
	}
}

func BuildTemplateCache() (map[string]*template.Template, error) {
	tc := map[string]*template.Template{}

	// Find out all the pages ending with page.gohtml in the templates directory
	pages, pageErr := filepath.Glob("./templates/*.page.gohtml")

	if pageErr != nil {
		log.Println("cannot find file in the specified path")
		return nil, errors.New("cannot build cache")
	}

	// Populate the cache
	for _, page := range pages {
		pageName := filepath.Base(page)
		var parseErr error
		tc[pageName], parseErr = template.ParseFiles("./templates/"+pageName, "./templates/base.layout.gohtml")

		if parseErr != nil {
			log.Println("Parse Error: failed to parse ", pageName)
		}
	}

	// Return the cache
	return tc, nil
}
