package postgres

import (
	"context"

	"fleet-management/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(
	db *pgxpool.Pool,
) *LocationRepository {
	return &LocationRepository{
		db: db,
	}
}

func (r *LocationRepository) Insert(
	ctx context.Context,
	location *model.VehicleLocation,
) error {

	query := `
		INSERT INTO vehicle_locations
		(
			vehicle_id,
			latitude,
			longitude,
			timestamp
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		location.VehicleID,
		location.Latitude,
		location.Longitude,
		location.Timestamp,
	)

	return err
}

func (r *LocationRepository) GetLatestByVehicleID(
	ctx context.Context,
	vehicleID string,
) (*model.VehicleLocation, error) {

	query := `
		SELECT
			id,
			vehicle_id,
			latitude,
			longitude,
			timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var location model.VehicleLocation

	err := r.db.QueryRow(
		ctx,
		query,
		vehicleID,
	).Scan(
		&location.ID,
		&location.VehicleID,
		&location.Latitude,
		&location.Longitude,
		&location.Timestamp,
	)

	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (r *LocationRepository) GetHistory(
	ctx context.Context,
	vehicleID string,
	start int64,
	end int64,
) ([]model.VehicleLocation, error) {

	query := `
		SELECT
			id,
			vehicle_id,
			latitude,
			longitude,
			timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1 AND timestamp >= $2 AND timestamp <= $3
		ORDER BY timestamp DESC
	`

	rows, err := r.db.Query(ctx, query, vehicleID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.VehicleLocation
	for rows.Next() {
		var loc model.VehicleLocation
		if err := rows.Scan(
			&loc.ID,
			&loc.VehicleID,
			&loc.Latitude,
			&loc.Longitude,
			&loc.Timestamp,
		); err != nil {
			return nil, err
		}
		history = append(history, loc)
	}

	return history, nil
}