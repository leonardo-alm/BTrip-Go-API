package main

import (
	"backend/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

var users = []models.User{
	{
		ID:        1,
		FirstName: "leonardo",
		LastName:  "almeida",
		Email:     "leo@email.com",
		Password:  "senha",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
var trips = []models.Trip{
	{
		ID:         1,
		TripStatus: "confirmed",
		TripType:   "client meeting",
		TripFrom:   "Porto Alegre, POA",
		TripTo:     "Frankfurt, FRA",
		TripDate:   time.Now(),
		CreatedBy:  1,
	},
	{
		ID:         2,
		TripStatus: "cancelled",
		TripType:   "client meeting",
		TripFrom:   "Porto Alegre, POA",
		TripTo:     "Sydney, SYD",
		TripDate:   time.Now(),
		CreatedBy:  1,
	},
}

func (app *application) AllTrips(w http.ResponseWriter, r *http.Request) {
	_ = app.writeJSON(w, http.StatusOK, trips)
}

func (app *application) GetTrip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tripID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	trip, err := FindTripByID(trips, tripID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, trip)
}

func (app *application) InsertTrip(w http.ResponseWriter, r *http.Request) {
	var trip models.Trip

	err := app.readJSON(w, r, &trip)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	tripID := len(trips) + 1

	trip.ID = tripID
	trip.CreatedAt = time.Now()
	trip.UpdatedAt = time.Now()

	trips = append(trips, trip)

	resp := JSONResponse{
		Error:   false,
		Message: "trip successfully added",
		Data:    trip,
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	var updatedTrip models.Trip

	err := app.readJSON(w, r, &updatedTrip)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	id := chi.URLParam(r, "id")
	tripID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	index, err := FindTripIndexByID(trips, tripID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	trips[index].UpdatedAt = time.Now()
	trips[index].TripStatus = updatedTrip.TripStatus
	trips[index].TripType = updatedTrip.TripType
	trips[index].TripFrom = updatedTrip.TripFrom
	trips[index].TripTo = updatedTrip.TripTo
	trips[index].TripDate = updatedTrip.TripDate

	resp := JSONResponse{
		Error:   false,
		Message: "trip successfully updated",
		Data:    &trips[index],
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tripID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	index, err := FindTripIndexByID(trips, tripID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	deletedTrip := trips[index]
	trips = append(trips[:index], trips[index+1:]...)

	resp := JSONResponse{
		Error:   false,
		Message: "trip successfully deleted",
		Data:    deletedTrip,
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {

	var user models.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	userID := len(users) + 1
	hashedPassword, err := user.HashPassword(user.Password)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user.ID = userID
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	users = append(users, user)

	resp := JSONResponse{
		Error:   false,
		Message: "user successfully added",
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, resp)

}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := FindUserByEmail(users, requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	u := jwtUser{
		ID:    user.ID,
		Email: user.Email,
	}

	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := FindUserByID(users, userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:    user.ID,
				Email: user.Email,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPairs)

		}
	}
	app.errorJSON(w, errors.New("refresh cookie not found"), http.StatusUnauthorized)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}
