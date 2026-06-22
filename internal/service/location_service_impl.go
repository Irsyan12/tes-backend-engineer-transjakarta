package service

import (
	"context"
	"fleet-management/internal/model"
	"fleet-management/internal/repository"
)

type locationServiceImpl struct {
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) LocationService {
	return &locationServiceImpl{
		repo: repo,
	}
}

func (s *locationServiceImpl) GetLatestByVehicleID(ctx context.Context, vehicleID string) (*model.VehicleLocation, error) {
	return s.repo.GetLatestByVehicleID(ctx, vehicleID)
}

func (s *locationServiceImpl) GetHistory(ctx context.Context, vehicleID string, start int64, end int64) ([]model.VehicleLocation, error) {
	return s.repo.GetHistory(ctx, vehicleID, start, end)
}

func (s *locationServiceImpl) Save(ctx context.Context, vehicleID string, latitude float64, longitude float64, timestamp int64) error {
	loc := &model.VehicleLocation{
		VehicleID: vehicleID,
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: timestamp,
	}
	return s.repo.Insert(ctx, loc)
}
