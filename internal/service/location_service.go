package service

import (
	"context"
	"fleet-management/internal/model"
)

type LocationService interface {
	GetLatestByVehicleID(
		ctx context.Context,
		vehicleID string,
	) (*model.VehicleLocation, error)

	GetHistory(
		ctx context.Context,
		vehicleID string,
		start int64,
		end int64,
	) ([]model.VehicleLocation, error)

	Save(
		ctx context.Context,
		vehicleID string,
		latitude float64,
		longitude float64,
		timestamp int64,
	) error
}