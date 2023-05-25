package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ASHWIN776/learning-Go/internal/config"
	"github.com/ASHWIN776/learning-Go/internal/driver"
	"github.com/ASHWIN776/learning-Go/internal/handlers"
	"github.com/ASHWIN776/learning-Go/internal/helpers"
	"github.com/ASHWIN776/learning-Go/internal/models"
	"github.com/ASHWIN776/learning-Go/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	// Defer the connection close
	defer db.SQL.Close()
	// Defer the mail channel
	defer close(app.MailChan)

	// Listens in the background
	log.Println("Starting mail listener")
	ListenForMail()

	// emailMsg := models.MailData{
	// 	From:    "john.do@gmail.com",
	// 	To:      "me@here.com",
	// 	Subject: "Test Email",
	// 	Content: "<h1>You got a mail.</h1><p>Check it out -></p>",
	// }

	// // Sending the emailMsg throught the app.MailChan (to share it with the go routine that deals with sending emails)
	// app.MailChan <- emailMsg

	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Printf("Serving on port %s\n", portNumber)
	err = srv.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start the server")
		return
	}

}

func run() (*driver.DB, error) {
	// What can I put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use Template cache")
	dbHost := flag.String("dbhost", "localhost", "Database Host")
	dbName := flag.String("dbname", "", "Database Name")
	dbUser := flag.String("dbuser", "", "Database User")
	dbPass := flag.String("dbpass", "", "Database Password")
	dbPort := flag.String("port", "5432", "Database Port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings(disable, prefer, require)")

	flag.Parse()

	// Exit if dbname and dbuser are not specified
	if *dbName == "" || *dbUser == "" {
		log.Println("Didn't specify required flags")
		os.Exit(1)
	}

	app.InProduction = *inProduction
	app.UseCache = *useCache

	tc, err := render.BuildTemplateCache()

	if err != nil {
		return nil, errors.New("cannot build template cache")
	}

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// Info log (will be a place to send all the info logs - currently to stdout)
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// Error log (will be a place to send all the Error logs - currently to stdout)
	errorLog := log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	app.Session = session
	app.TemplateCache = tc

	// Creating a channel for mail data and assigning it to app.MailChan
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// Connecting to database
	connString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	log.Println(connString)
	db, err := driver.ConnectSQL(connString)

	if err != nil {
		return nil, err
	}
	log.Println("Connected to Database")

	render.GetConfig(&app)
	helpers.GetConfig(&app)

	// As Home and About are Methods of an instance of type Repository - So, we need the instance from render.go
	handlers.AddRepo(&app, db)

	return db, nil
}
