package main

import (
	"net/http"
	"testing"
)

func TestWriteToConsole(t *testing.T) {
	var myH http.Handler

	h := WriteToConsole(myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("does not return a http.Handler, returns a %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("does not return a http.Handler, returns a %T", v)
	}
}

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := Nosurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("does not return a http.Handler, returns a %T", v)
	}
}
