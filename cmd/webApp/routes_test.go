package main

import (
	"testing"

	"github.com/ASHWIN776/learning-Go/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {

	mux := routes(&config.AppConfig{})

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("does not return a *chi.Mux, returns a %T", v)
	}
}
