package forms

import (
	"net/http"
	"net/url"
	"testing"
)

type postData struct {
	field string
	value string
}

func TestForm_Valid(t *testing.T) {
	r, err := http.NewRequest("POST", "/some", nil)

	if err != nil {
		t.Error("could not create request")
	}

	form := New(r.PostForm)

	isValid := form.IsValid()

	if !isValid {
		t.Error("Failed IsValid()")
	}
}

var hasTests = []struct {
	formData      postData
	field         string
	expectedValue bool
}{
	{
		formData:      postData{"name", "John"},
		field:         "name",
		expectedValue: true,
	},
	{
		formData:      postData{"name", ""},
		field:         "name",
		expectedValue: false,
	},
	{
		formData:      postData{},
		field:         "name",
		expectedValue: false,
	},
	{
		formData:      postData{"name", "John"},
		field:         "email",
		expectedValue: false,
	},
}

func TestForm_Has(t *testing.T) {
	for _, test := range hasTests {

		formPostedData := url.Values{}
		formPostedData.Add(test.formData.field, test.formData.value)

		form := New(formPostedData)

		isValid := form.Has(test.field)

		if isValid != test.expectedValue {
			t.Errorf("expected value was %t, got %t", test.expectedValue, isValid)
		}
	}
}

func TestForm_Required(t *testing.T) {
	r, err := http.NewRequest("POST", "/some", nil)

	if err != nil {
		t.Error("could not create request")
	}

	form := New(r.PostForm)
	form.Required("name")
	isValid := form.IsValid()

	if isValid {
		t.Error("shows valid when given field does not exist")
	}

	postData := url.Values{}
	postData.Add("name", "John Doe")
	postData.Add("email", "j@gmail.com")

	form = New(postData)
	form.Required("name", "email")
	isValid = form.IsValid()

	if !isValid {
		t.Error("shows invalid when given fields and their values exist")
	}

}

func TestForm_MinLength(t *testing.T) {
	postData := url.Values{}
	postData.Add("name", "John")

	form := New(postData)

	form.MinLength("name", 5)
	isValid := form.IsValid()

	if isValid {
		t.Error("field value doesnot satisfy the minlength and still returns valid")
	}

	postData.Set("name", "Ashwin")
	form = New(postData) // need a new form cuz otherwise ill have to remove the previous error
	form.MinLength("name", 4)
	isValid = form.IsValid()

	if !isValid {
		t.Error("field value satisfies minlength but returns invalid")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	postData.Add("email", "john")

	form := New(postData)
	form.IsEmail("email")
	isValid := form.IsValid()

	if isValid {
		t.Error("shows valid when the email is invalid")
	}

	postData.Set("email", "a@gmail.com")
	form = New(postData)
	form.IsEmail("email")
	isValid = form.IsValid()

	if !isValid {
		t.Error("shows invalid when the email is valid")
	}
}
