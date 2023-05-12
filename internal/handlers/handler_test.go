package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var tests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
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
			t.Logf("passed %s", test.name)
		} else {
			// For POST reqs
		}
	}
}
