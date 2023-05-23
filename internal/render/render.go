package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ASHWIN776/learning-Go/internal/config"
	"github.com/ASHWIN776/learning-Go/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"

// Gets the pointer to the app config created in main.go, and assigns it to the pointer of the same type created here
func GetConfig(a *config.AppConfig) {
	app = a
}

func addDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1 // otherwise it will put the zero value of int there(0)
	}

	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, file string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// Get the template cache from the App config
		tc = app.TemplateCache
	} else {
		tc, _ = BuildTemplateCache()
	}

	td = addDefaultData(td, r)

	execErr := tc[file].Execute(w, td)

	if execErr != nil {
		log.Println(execErr)
		return execErr
	}

	return nil
}

func BuildTemplateCache() (map[string]*template.Template, error) {

	tc := map[string]*template.Template{}

	// Find out all the pages ending with page.gohtml in the templates directory
	pages, pageErr := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))

	if pageErr != nil {
		log.Println("cannot find file in the specified path")
		return nil, errors.New("cannot build cache")
	}

	// Find out all the pages ending with layout.gohtml in the template directory
	layoutPages, pageErr := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
	if pageErr != nil {
		log.Println("cannot find layout file in the specified path")
		return nil, errors.New("cannot build cache")
	}

	// Populate the cache
	for _, page := range pages {
		pageName := filepath.Base(page)
		var parseErr error

		// parse the individual file
		cache, parseErr := template.ParseFiles(fmt.Sprintf("%s/%s", pathToTemplates, pageName))

		if parseErr != nil {
			log.Println("Parse Error: failed to parse ", pageName)
			return nil, parseErr
		}

		// parse the layout files too
		if len(layoutPages) > 0 {
			cache, parseErr = cache.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if parseErr != nil {
				log.Println("Parse Error: failed to parse ", pageName)
				return nil, parseErr
			}
		}

		// Put the cache in tc
		tc[pageName] = cache
	}

	// Return the cache
	return tc, nil
}
