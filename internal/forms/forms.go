package forms

import (
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	Data   url.Values
	Errors errors
}

// Initialises a Form struct
func New(data url.Values) *Form {
	return &Form{
		Data:   data,
		Errors: errors{},
	}
}

// Checks for validitiy, if errors are present returns false
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}

// Checks for a reqd field (if it has as value or not)
func (f *Form) Has(field string, r *http.Request) bool {
	val := r.Form.Get(field)

	if val == "" {
		f.Errors.Add("firstName", "cannot leave this field blank")
		return false
	}

	return true
}

// Checks for required fields
func (f *Form) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		val := r.Form.Get(string(field))

		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, "This field is required")
		}
	}
}
