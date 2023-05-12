package render

import (
	"log"
	"net/http"
	"testing"

	"github.com/ASHWIN776/learning-Go/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td = &models.TemplateData{}
	r, err := CreateRequest()

	if err != nil {
		t.Error("cannot create request")
	}

	session.Put(r.Context(), "flash", "worked")

	result := addDefaultData(td, r)

	if result.Flash != "worked" {
		t.Errorf("expected to get worked, got %s", result.Flash)
	}

}

func CreateRequest() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)

	if err != nil {
		log.Fatal(err)
	}

	// Adding session into the request
	/*
		1. Create a context
		2. Put session data in there
		3. Put the context back into the request made above
	*/
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))

	r = r.WithContext(ctx)

	return r, nil
}
