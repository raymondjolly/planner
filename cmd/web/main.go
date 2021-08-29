package main

import (
	"fmt"
	"github.com/raymondjolly/bookings/internal/config"
	"github.com/raymondjolly/bookings/internal/handlers"
	"github.com/raymondjolly/bookings/internal/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portValue string = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = time.Hour * 24
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	fatalErrCheck(err)

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("The web server has started and is available on port: %s\n", portValue)
	// http.ListenAndServe(portValue, nil)
	srv := &http.Server{
		Addr:    portValue,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	fatalErrCheck(err)
}

func fatalErrCheck(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
