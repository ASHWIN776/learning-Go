package models

import "github.com/ASHWIN776/learning-Go/internal/forms"

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float64
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form // to display errors after a user submits the form and a server side validation error occurs
	IsAuthenticated int         // If greater than 0, user is authenticated, == 0, then not
}
