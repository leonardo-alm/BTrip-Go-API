package models

import "time"

type Trip struct {
	ID         int       `json:"id"`
	TripStatus string    `json:"tripStatus"`
	TripType   string    `json:"tripType"`
	TripFrom   string    `json:"tripFrom"`
	TripTo     string    `json:"tripTo"`
	TripDate   time.Time `json:"tripDate"`
	CreatedBy  int       `json:"createdBy"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
