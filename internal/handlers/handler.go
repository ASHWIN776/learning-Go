package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ASHWIN776/learning-Go/internal/config"
	"github.com/ASHWIN776/learning-Go/internal/driver"
	"github.com/ASHWIN776/learning-Go/internal/forms"
	"github.com/ASHWIN776/learning-Go/internal/helpers"
	"github.com/ASHWIN776/learning-Go/internal/models"
	"github.com/ASHWIN776/learning-Go/internal/render"
	"github.com/ASHWIN776/learning-Go/internal/repository"
	"github.com/ASHWIN776/learning-Go/internal/repository/dbrepo"
)

type Repository struct {
	app *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo Repository

func AddRepo(a *config.AppConfig, db *driver.DB) {
	Repo.app = a
	Repo.DB = dbrepo.NewPostgresRepo(db.SQL, a)
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
	query := r.URL.Query().Get("error")
	error := ""
	if query == "true" {
		error = "create a reservation first"
	}

	data := make(map[string]interface{})
	data["resDetails"] = models.Reservation{} // Creating an empty res

	// Render template
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form:  forms.New(nil),
		Data:  data,
		Error: error,
	})
}

func (rep *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Mon Jan 2 15:04:05 MST 2006  (MST is GMT-0700) - 01/02 03:04:05PM '06 -0700
	layout := "02-01-2006"

	sd := r.Form.Get("startDate")
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}

	ed := r.Form.Get("endDate")
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomId, err := strconv.Atoi(r.Form.Get("roomId"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	// I can send this to the form to re-render if there are any errors
	resDetails := models.Reservation{
		FirstName: r.Form.Get("firstName"),
		LastName:  r.Form.Get("lastName"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phoneNumber"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomId,
	}

	form := forms.New(r.PostForm)

	// form.Has("firstName", r)
	form.Required("firstName", "lastName", "email")
	form.MinLength("firstName", 3)
	form.IsEmail("email")

	isValid := form.IsValid()

	data := make(map[string]interface{})
	data["resDetails"] = resDetails

	if !isValid {

		render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	// Insert reservation in the database(reservations table)
	resId, err := rep.DB.InsertReservation(resDetails)
	if err != nil {
		helpers.ServerError(w, err)
	}

	// Insert this into the room_restrictions table
	log.Println("Reservation Id: ", resId)

	rep.app.Session.Put(r.Context(), "resDetails", resDetails)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

func (rep *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// Perform some logic
	resDetails, ok := rep.app.Session.Get(r.Context(), "resDetails").(models.Reservation)

	if !ok {
		rep.app.Session.Put(r.Context(), "error", "no reservation found")
		rep.app.ErrorLog.Println("no reservation found")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["resDetails"] = resDetails

	// Render Template
	render.RenderTemplate(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

func (rep *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {

	// Perform some logic

	// Render template
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{
		Form: forms.New(nil), // to display the errors and other info after a user submits
	})
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
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
}
