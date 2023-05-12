package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"POST search-availability-json", "/search-availability-json", "POST", []postData{
		{"startDate", "01-01-2023"},
		{"endDate", "02-01-2023"},
	}, http.StatusOK},
	{"POST search-availability", "/search-availability", "POST", []postData{
		{"startDate", "01-01-2023"},
		{"endDate", "02-01-2023"},
	}, http.StatusOK},
	{"POST make-reservation", "/make-reservation", "POST", []postData{
		{"firstName", "Ashwin"},
		{"lastName", "Anil"},
		{"email", "a@gmail.com"},
		{"phoneNumber", "1234-2345-444"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := GetRoutes()

	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range tests {
		if test.method == "GET" {
			// For GET reqs

			// Make a req as a client
			res, err := testServer.Client().Get(testServer.URL + test.url)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != test.expectedStatusCode {
				t.Errorf("expected status code for %s is %d, got %d", test.name, test.expectedStatusCode, res.StatusCode)
			}
		} else {
			// For POST reqs
			values := url.Values{}

			for _, param := range test.params {
				values.Add(param.key, param.value)
			}

			res, err := testServer.Client().PostForm(testServer.URL+test.url, values)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != test.expectedStatusCode {
				t.Errorf("expected status code for %s is %d, got %d", test.name, test.expectedStatusCode, res.StatusCode)
			}
		}
	}
}
