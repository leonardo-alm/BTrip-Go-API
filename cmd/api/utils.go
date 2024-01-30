package main

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

func FindTripByID(trips []models.Trip, id int) (*models.Trip, error) {
	for _, trip := range trips {
		if trip.ID == id {
			return &trip, nil
		}
	}

	return nil, fmt.Errorf("Trip with ID %d not found", id)
}

func FindTripIndexByID(trips []models.Trip, id int) (int, error) {
	for i, trip := range trips {
		if trip.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("Trip with ID %d not found", id)
}

func FindUserByEmail(users []models.User, email string) (*models.User, error) {
	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("email not found")
}

func FindUserByID(users []models.User, id int) (*models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user with ID %d not found", id)
}
