package handlers

import (
	"myApp/pkg/config"
	"myApp/pkg/models"
	"myApp/pkg/render"
	"net/http"
)

type Repository struct {
	app *config.AppConfig
}

var Repo Repository

func GetRepo(a *config.AppConfig) *Repository {
	Repo.app = a

	return &Repo
}

func (rep *Repository) Home(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["text"] = "Hello, this is Ashwin Anil"

	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{})
}
