package main

import (
	"net/http"

	"github.com/imsujan276/golang-bookings/internal/config"
	"github.com/imsujan276/golang-bookings/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// custom middleware
	r.Use(NoSurf)
	r.Use(SessionLoad)

	r.Get("/", handlers.Repo.Home)
	r.Get("/about", handlers.Repo.About)
	r.Get("/generals-quarters", handlers.Repo.Generals)
	r.Get("/majors-suite", handlers.Repo.Majors)
	r.Get("/contact", handlers.Repo.Contact)

	r.Get("/make-reservation", handlers.Repo.Reservation)
	r.Post("/make-reservation", handlers.Repo.PostReservation)
	r.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	r.Get("/search-availability", handlers.Repo.SearchAvailability)
	r.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	r.Post("/search-availability-json", handlers.Repo.JsonSearchAvailability)

	r.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	r.Get("/book-room", handlers.Repo.BookRoom)

	r.Get("/login", handlers.Repo.ShowLogin)
	r.Post("/login", handlers.Repo.PostShowLogin)
	r.Get("/logout", handlers.Repo.Logout)

	// Admin/Authenticated route
	r.Route("/admin", func(r chi.Router) {
		r.Use(Auth)
		r.Get("/dashboard", handlers.Repo.AdminDashboard)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return r

}
