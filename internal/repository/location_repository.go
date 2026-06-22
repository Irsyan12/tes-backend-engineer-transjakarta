package repository

import (
	"context"
	"fleet-management/internal/model"
)

type LocationRepository interface {
	Insert(ctx context.Context, location *model.VehicleLocation) error

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
}