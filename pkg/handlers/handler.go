package handlers

import (
	"myApp/pkg/config"
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
	render.RenderTemplate(w, "home.page.gohtml")
}

func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.gohtml")
}
