package main

import (
	"fmt"
	"net/http"

	"github.com/ASHWIN776/learning-Go/internal/helpers"
	"github.com/justinas/nosurf"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Nosurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func Auth(next http.Handler) http.Handler {
	// The below HandlerFunc() will give the passed in function access to the response writer and the request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			// Put an error alert msg in the session and redirect to the login screen
			session.Put(r.Context(), "error", "user not logged in")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// for running the next middleware or the handler itself
		next.ServeHTTP(w, r)
	})
}
