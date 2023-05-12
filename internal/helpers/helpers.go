package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ASHWIN776/learning-Go/internal/config"
)

var app *config.AppConfig

func GetConfig(a *config.AppConfig) {
	app = a
}

// status - code of error given by the server error
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("client error with status of ", status)

	// Prints it on the web
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	// Useful info to show as error
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.ErrorLog.Println(trace)

	// Prints it on the web
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
