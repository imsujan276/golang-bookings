package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/imsujan276/golang-bookings/internal/config"
	"github.com/imsujan276/golang-bookings/internal/handlers"
	"github.com/imsujan276/golang-bookings/internal/models"
	"github.com/imsujan276/golang-bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// define the type of data that can be stored into session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	// session lasts for 24 hours
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // use ssl; set to true in production

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("Can not create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

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
