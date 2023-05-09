package handlers

import (
	"encoding/json"
	"fmt"
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

	// Perform some logic
	remoteIp := r.RemoteAddr
	rep.app.Session.Put(r.Context(), "remote_ip", remoteIp)

	stringMap := make(map[string]string)
	stringMap["text"] = "Hello, this is Ashwin Anil"

	// Render template
	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {

	// Perform some logic
	stringMap := make(map[string]string)
	stringMap["remote_ip"] = rep.app.Session.GetString(r.Context(), "remote_ip")

	// Render template
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (rep *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	// Perform some logic

	// Render template
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}

func (rep *Repository) GeneralsQuarters(w http.ResponseWriter, r *http.Request) {

	// Perform some logic

	// Render template
	render.RenderTemplate(w, r, "generals-quarters.page.gohtml", &models.TemplateData{})
}

func (rep *Repository) MajorsSuite(w http.ResponseWriter, r *http.Request) {

	// Perform some logic

	// Render template
	render.RenderTemplate(w, r, "majors-suite.page.gohtml", &models.TemplateData{})
}

func (rep *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {

	// Perform some logic

	// Render template
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{})
}

func (rep *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {

	// Perform some logic

	// Render template
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

func (rep *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	startDate := r.Form.Get("startDate")
	endDate := r.Form.Get("endDate")

	w.Write([]byte(fmt.Sprintf("Start date is %s and End date is %s", startDate, endDate)))
}

type JSONResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (rep *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	res := JSONResponse{
		Ok:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(res, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
}
