package handlers

import (
	"net/http"

	"github.com/ASHWIN776/learning-Go/pkg/config"
	"github.com/ASHWIN776/learning-Go/pkg/models"
	"github.com/ASHWIN776/learning-Go/pkg/render"
)

type Repository struct {
	app *config.AppConfig
}

var Repo Repository

func AddRepo(a *config.AppConfig) {
	Repo.app = a
}

func (rep *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIp := r.RemoteAddr
	rep.app.Session.Put(r.Context(), "remote_ip", remoteIp)

	stringMap := make(map[string]string)
	stringMap["text"] = "Hello, this is Ashwin Anil"

	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["remote_ip"] = rep.app.Session.GetString(r.Context(), "remote_ip")

	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}
