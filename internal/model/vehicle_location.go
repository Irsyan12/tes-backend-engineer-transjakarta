package model

import "time"

type VehicleLocation struct {
	ID        int64     `json:"id"`
	VehicleID string    `json:"vehicle_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp int64     `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}
