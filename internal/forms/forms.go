package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
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

// Checks if a field satisfies the given length
func (f *Form) MinLength(field string, length int, r *http.Request) {
	val := r.Form.Get(field)

	if len(val) < length {
		f.Errors.Add(field, fmt.Sprintf("given value should be atleast %d letters", length))
	}
}

// Checks if a field has valid email as its value
func (f *Form) IsEmail(field string, r *http.Request) {
	val := r.Form.Get(field)

	if !govalidator.IsEmail(val) {
		f.Errors.Add(field, "not a valid email")
	}
}
