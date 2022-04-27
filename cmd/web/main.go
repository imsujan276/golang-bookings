package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/imsujan276/golang-bookings/internal/config"
	"github.com/imsujan276/golang-bookings/internal/driver"
	"github.com/imsujan276/golang-bookings/internal/handlers"
	"github.com/imsujan276/golang-bookings/internal/helpers"
	"github.com/imsujan276/golang-bookings/internal/models"
	"github.com/imsujan276/golang-bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChannel)
	listenForMail()

	fmt.Println("Serving application on port", portNumber)
	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {
	// define the type of data that can be stored into session
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(models.Reservation{})

	mailChannel := make(chan models.MailData)
	app.MailChannel = mailChannel

	// change this to true when in production
	app.InProduction = false

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Llongfile)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	// session lasts for 24 hours
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // use ssl; set to true in production

	app.Session = session

	//connect to database
	log.Println("Connecting to Database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=password123")
	if err != nil {
		log.Fatal("Can not connect to database! Dying...")
	}
	log.Println("Database Connected...")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("Can not create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
