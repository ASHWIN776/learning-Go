package handlers

import (
	"encoding/json"
	"errors"
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
	"github.com/go-chi/chi/v5"
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

	data := make(map[string]interface{})
	res, ok := rep.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get from session"))
		return
	}
	// Get the room info - using the roomId and put it in res, and update the session --------
	room, err := rep.DB.GetRoomById(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room = room

	rep.app.Session.Put(r.Context(), "reservation", res) // Updated the session
	// ---------------------------------------------------------------

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	var stringMap = make(map[string]string)
	stringMap["startDate"] = sd
	stringMap["endDate"] = ed

	data["reservation"] = res

	// Render template
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (rep *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	reservation, ok := rep.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// I can send this to the form to re-render if there are any errors
	reservation.FirstName = r.Form.Get("firstName")
	reservation.LastName = r.Form.Get("lastName")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phoneNumber")

	form := forms.New(r.PostForm)

	// form.Has("firstName", r)
	form.Required("firstName", "lastName", "email")
	form.MinLength("firstName", 3)
	form.IsEmail("email")

	isValid := form.IsValid()

	data := make(map[string]interface{})
	data["reservation"] = reservation

	if !isValid {
		var stringMap = make(map[string]string)
		stringMap["startDate"] = reservation.StartDate.Format("2006-01-02")
		stringMap["endDate"] = reservation.EndDate.Format("2006-01-02")

		render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})

		return
	}

	// Insert reservation in the database(reservations table)
	resId, err := rep.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Insert this into the room_restrictions table
	restrictionDetails := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: resId,
		RestrictionID: 1,
	}

	err = rep.DB.InsertRoomRestriction(restrictionDetails)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Send notification email - first to guest------------------
	htmlMessage := fmt.Sprintf(`
		<h1>Booking Confirmation</h1>
		<span>Hey %s</span><br>
		<span>Your booking from %s to %s is confirmed</span>
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	emailMsg := models.MailData{
		From:     "john.do@gmail.com",
		To:       reservation.Email,
		Subject:  "Booked! Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	// Sending the emailMsg throught the app.MailChan (to share it with the go routine that deals with sending emails)
	rep.app.MailChan <- emailMsg
	// ----------------------------------------------------------

	// Send notification email - to the property owner (john.do@gmail.com)
	htmlMessage = fmt.Sprintf(`
		<h1>Booking Notification</h1>
		<span>Hey John</span><br>
		<span>Booking added by %s from %s to %s.</span>
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	emailMsg = models.MailData{
		From:     "john.do@gmail.com",
		To:       "john.do@gmail.com",
		Subject:  "+1 Booking Notification",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	rep.app.MailChan <- emailMsg
	// ----------------------------------------------------------

	// Updating the session again
	rep.app.Session.Put(r.Context(), "reservation", reservation)

	// Redirecting to the Reservation Summary Page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

func (rep *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// Perform some logic
	reservation, ok := rep.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rep.app.Session.Put(r.Context(), "error", "no reservation found")
		rep.app.ErrorLog.Println("no reservation found")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	stringMap := make(map[string]string)
	stringMap["startDate"] = reservation.StartDate.Format("2006-01-02")
	stringMap["endDate"] = reservation.EndDate.Format("2006-01-02")

	// Render Template
	render.RenderTemplate(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
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
	// Parsing the form to populate r.Form map
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	layout := "2006-01-02"

	sd := r.Form.Get("startDate")
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	ed := r.Form.Get("endDate")
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get available rooms
	rooms, err := rep.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// If not available - error message in session
	if len(rooms) == 0 {
		rep.app.Session.Put(r.Context(), "error", "no available rooms")
		log.Println("No rooms")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	// Show rooms in template(if available)
	// Inserting the reservation details(startDate and endDate) into the session
	res := &models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	rep.app.Session.Put(r.Context(), "reservation", res)

	// Sending rooms to the template
	var data = make(map[string]interface{})
	data["rooms"] = rooms

	// Render the choose-room.page.gohtml page
	render.RenderTemplate(w, r, "choose-room.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

type JSONResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"roomId"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// Gets the dates and roomId and sends the response containing the room's availability
func (rep *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	// Get the required info from the form - roomId, and dates(to call the SearchAvailabilityByRoomID)
	roomId, err := strconv.Atoi(r.Form.Get("roomId"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("startDate")
	startDate, err := time.Parse("2006-01-02", sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	ed := r.Form.Get("endDate")
	endDate, err := time.Parse("2006-01-02", ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	isAvailable, err := rep.DB.SearchAvailabilityByRoomId(startDate, endDate, roomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Sending the availability as JSON response
	res := JSONResponse{
		Ok:        isAvailable,
		Message:   "",
		RoomID:    strconv.Itoa(roomId),
		StartDate: sd,
		EndDate:   ed,
	}

	out, err := json.MarshalIndent(res, "", "    ")

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
}

// Gets the roomId from the query param after the user clicks on the required room(available) and redirects the user to /make-reservation after updating the session with the roomId
func (rep *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// Get the roomId from the url param
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get the reservation detail from session and add room id
	res, ok := rep.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get from session"))
	}

	res.RoomID = roomId

	// Put the reservation detail back into the session and redirect the page to make-reservation
	rep.app.Session.Put(r.Context(), "reservation", res)

	// Redirect to /make-reservation
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (rep *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("roomId"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	startDate, err := time.Parse("2006-01-02", r.URL.Query().Get("s"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse("2006-01-02", r.URL.Query().Get("e"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Creating a reservation - to be put in the session and passed on to /make-reservation
	reservation := models.Reservation{
		RoomID:    id,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// reservation added to the session
	rep.app.Session.Put(r.Context(), "reservation", reservation)

	// Redirect to /make-reservation
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (rep *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {

	// Render the login template
	render.RenderTemplate(w, r, "login.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Handler associated with logging the user in
func (rep *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	// Renew Session Token to prevent Session Fixation Attack

	// Get the values from the form
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, errors.New("cannot parse the form"))
	}

	// Server Side Validation ---------------------------------------------
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	form.MinLength("password", 8)

	isValid := form.IsValid()

	if !isValid {
		var data = make(map[string]interface{})
		data["email"] = r.Form.Get("email")

		render.RenderTemplate(w, r, "login.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}
	// --------------------------------------------------------------------

	// Authenticate the creds entered
	enteredEmail := r.Form.Get("email")
	enteredPass := r.Form.Get("password")
	userId, _, err := rep.DB.Authenticate(enteredEmail, enteredPass)
	if err != nil {
		rep.app.ErrorLog.Println(err)

		// Put an error message in the session for the alert, and redirect back to the login page
		rep.app.Session.Put(r.Context(), "error", "Login failed: Invalid user credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Put the user id in the session
	rep.app.Session.Put(r.Context(), "user_id", userId)

	// Redirect to the home, with the success alert
	rep.app.Session.Put(r.Context(), "flash", "Logged In Successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// logs a user out
func (rep *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	// Destroy the session - Could have done the .Remove on the key(user_id), but who knows what all key i will add in the future which would contain the logged in user detail. I will have to write the .Remove for all of those keys
	_ = rep.app.Session.Destroy(r.Context())

	// Renew Token
	_ = rep.app.Session.RenewToken(r.Context())

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (rep *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "admin-dashboard.page.gohtml", &models.TemplateData{})
}

// Shows all the new reservations in the database(the ones that are not processed)
func (rep *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {

	newReservations, err := rep.DB.NewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["newReservations"] = newReservations

	render.RenderTemplate(w, r, "admin-new-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

// Shows all the reservations in the database
func (rep *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {

	allReservations, err := rep.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["allReservations"] = allReservations

	render.RenderTemplate(w, r, "admin-all-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

// Shows more of one reservation
func (rep *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	pageSrc := chi.URLParam(r, "src")
	reservationId, err := strconv.Atoi(chi.URLParam(r, "reservation_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation, err := rep.DB.GetReservationById(reservationId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	stringMap := make(map[string]string)
	stringMap["src"] = pageSrc

	render.RenderTemplate(w, r, "admin-show-reservation.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

// Processes the form values when saved(show reservation)
func (rep *Repository) PostAdminShowReservation(w http.ResponseWriter, r *http.Request) {
	pageSrc := chi.URLParam(r, "src")
	reservationId, err := strconv.Atoi(chi.URLParam(r, "reservation_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get the old reservation using the reservationId
	oldReservation, err := rep.DB.GetReservationById(reservationId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Server Side validation
	form := forms.New(r.PostForm)

	form.Required("firstName", "lastName", "email")
	form.IsEmail("email")

	if !form.IsValid() {
		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
		return
	}

	// Update oldReservation with values submitted from the form
	oldReservation.FirstName = r.Form.Get("firstName")
	oldReservation.LastName = r.Form.Get("lastName")
	oldReservation.Email = r.Form.Get("email")
	oldReservation.Phone = r.Form.Get("phoneNumber")

	// Update the reservations table
	err = rep.DB.UpdateReservation(oldReservation) // oldReservation contains the new reservation values
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rep.app.Session.Put(r.Context(), "flash", fmt.Sprintf("Reservation ID: %d Updated Successfully", oldReservation.ID))
	http.Redirect(w, r, "/admin/reservations-"+pageSrc, http.StatusSeeOther)
}

// Processes the given reservation
func (rep *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	pageSrc := chi.URLParam(r, "src")
	reservationId, err := strconv.Atoi(chi.URLParam(r, "reservation_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Update the processed value in the reservation table row for the particular reservationId
	err = rep.DB.UpdateProcessedForReservation(reservationId, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// add the flash message in the session too
	rep.app.Session.Put(r.Context(), "flash", fmt.Sprintf("Processed Reservation ID %d successfully", reservationId))

	// Redirect the user back to the src page
	http.Redirect(w, r, "/admin/reservations-"+pageSrc, http.StatusSeeOther)
}

// Delete the reservation
func (rep *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	pageSrc := chi.URLParam(r, "src")
	reservationId, err := strconv.Atoi(chi.URLParam(r, "reservation_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Update the processed value in the reservation table row for the particular reservationId
	err = rep.DB.DeleteReservation(reservationId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// add the flash message in the session too
	rep.app.Session.Put(r.Context(), "flash", fmt.Sprintf("Deleted Reservation ID %d successfully", reservationId))

	// Redirect the user back to the src page
	http.Redirect(w, r, "/admin/reservations-"+pageSrc, http.StatusSeeOther)
}

// Shows reservation calendar
func (rep *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// if month and year are not specified
	now := time.Now()

	// if month and year are specified
	if r.URL.Query().Get("y") != "" {
		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	previous := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	previousMonth := previous.Format("01")
	previousMonthYear := previous.Format("2006")

	thisMonth := now.Format("01")
	thisMonthYear := now.Format("2006")

	// putting all these details in a stringmap, so as to be passed down to the template
	stringMap := make(map[string]string)
	stringMap["nextMonth"] = nextMonth
	stringMap["nextMonthYear"] = nextMonthYear
	stringMap["previousMonth"] = previousMonth
	stringMap["previousMonthYear"] = previousMonthYear
	stringMap["thisMonth"] = thisMonth
	stringMap["thisMonthYear"] = thisMonthYear

	// Get the number of days in this month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	days := lastOfMonth.Day()

	intMap := make(map[string]int)
	intMap["daysInMonth"] = days

	// Data map
	data := make(map[string]interface{})
	data["now"] = now

	// Get rooms info and pass into the data map
	rooms, err := rep.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	render.RenderTemplate(w, r, "admin-reservations-calendar.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}
