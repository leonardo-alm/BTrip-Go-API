package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Post("/register", app.register)
	mux.Post("/login", app.login)
	mux.Get("/logout", app.logout)

	mux.Route("/api", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/trips", app.AllTrips)
		mux.Post("/trips", app.InsertTrip)
		mux.Get("/trips/{id}", app.GetTrip)
		mux.Patch("/trips/{id}", app.UpdateTrip)
		mux.Delete("/trips/{id}", app.DeleteTrip)

		mux.Get("/users", app.AllUsers)

		mux.Get("/refresh", app.refreshToken)

	})

	return mux
}
