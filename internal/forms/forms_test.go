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

var requiredTests = []struct {
	formData      []postData
	fields        []string
	expectedValue bool
}{
	{
		formData:      []postData{},
		fields:        []string{"name", "email"},
		expectedValue: false,
	},
	{
		formData: []postData{
			{"name", "John"},
			{"email", "a@gmail.com"},
		},
		fields:        []string{"name", "email"},
		expectedValue: true,
	},
	{
		formData: []postData{
			{"name", ""},
			{"email", "a@gmail.com"},
		},
		fields:        []string{"name", "email"},
		expectedValue: false,
	},
	{
		formData: []postData{
			{"name", "John"},
			{"email", "a@gmail.com"},
		},
		fields:        []string{"name", "phoneNumber"},
		expectedValue: false,
	},
}

func TestForm_Required(t *testing.T) {

	for _, test := range requiredTests {
		formPostData := url.Values{}

		for _, data := range test.formData {
			formPostData.Add(data.field, data.value)
		}

		form := New(formPostData)
		form.Required(test.fields...)
		isValid := form.IsValid()

		if isValid != test.expectedValue {
			t.Errorf("expected %t, got %t for %v postData and %v fields", test.expectedValue, isValid, test.formData, test.fields)
		}
	}
}

var minLengthTests = []struct {
	formData      postData
	minLength     int
	expectedValue bool
}{
	{
		formData:      postData{"name", "John"},
		minLength:     3,
		expectedValue: true,
	},
	{
		formData:      postData{"name", "John"},
		minLength:     5,
		expectedValue: false,
	},
	{
		formData:      postData{"name", ""},
		minLength:     5,
		expectedValue: false,
	},
	{
		formData:      postData{},
		minLength:     5,
		expectedValue: false,
	},
}

func TestForm_MinLength(t *testing.T) {
	for _, test := range minLengthTests {
		formPostData := url.Values{}
		formPostData.Add(test.formData.field, test.formData.value)

		form := New(formPostData)
		form.MinLength(test.formData.field, test.minLength)
		isValid := form.IsValid()

		if isValid != test.expectedValue {
			t.Errorf("expected length is %d, given value(%s) length is %d", test.minLength, test.formData.value, len(test.formData.value))
		}

		// Checking errors.go (Get function)
		err := form.Errors.Get(test.formData.field)

		if isValid {
			if err != "" {
				t.Error("got error when not expecting one")
			}
		}

		if !isValid {
			if err == "" {
				t.Error("expected error, did not get it")
			}
		}
	}
}

var isEmailTests = []struct {
	formData      postData
	expectedValue bool
}{
	{
		formData:      postData{"email", "a@gmail.com"},
		expectedValue: true,
	},
	{
		formData:      postData{"email", "John"},
		expectedValue: false,
	},
	{
		formData:      postData{"email", ""},
		expectedValue: false,
	},
	{
		formData:      postData{},
		expectedValue: false,
	},
}

func TestForm_IsEmail(t *testing.T) {
	for _, test := range isEmailTests {
		formPostData := url.Values{}
		formPostData.Add(test.formData.field, test.formData.value)

		form := New(formPostData)
		form.IsEmail(test.formData.field)
		isValid := form.IsValid()

		if isValid != test.expectedValue {
			t.Errorf("expected %t, got %t for email value: %s", test.expectedValue, isValid, test.formData.value)
		}
	}
}
