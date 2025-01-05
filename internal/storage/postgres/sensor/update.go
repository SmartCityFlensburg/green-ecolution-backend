package sensor

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (r *SensorRepository) Update(ctx context.Context, id string, sFn ...entities.EntityFunc[entities.Sensor]) (*entities.Sensor, error) {
	entity, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, fn := range sFn {
		fn(entity)
	}

	if err := r.updateEntity(ctx, entity); err != nil {
		slog.Error("failed to update sensor entity", "Error", err)
		return nil, err
	}

	if entity.LatestData != nil && entity.LatestData.Data != nil {
		err = r.InsertSensorData(ctx, entity.LatestData, entity.ID)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, entity.ID)
}

func (r *SensorRepository) updateEntity(ctx context.Context, sensor *entities.Sensor) error {
	params := sqlc.UpdateSensorParams{
		ID:     sensor.ID,
		Status: sqlc.SensorStatus(sensor.Status),
	}

	locationParams := &sqlc.SetSensorLocationParams{
		ID:        sensor.ID,
		Latitude:  sensor.Latitude,
		Longitude: sensor.Longitude,
	}

	if err := r.validateCoordinates(locationParams); err != nil {
		return err
	}
	err := r.store.SetSensorLocation(ctx, locationParams)
	if err != nil {
		slog.Error("failed to update sensor location", "Error", err)
		return err
	}

	return r.store.UpdateSensor(ctx, &params)
}
func (r *SensorRepository) validateCoordinates(locationParams *sqlc.SetSensorLocationParams) error {
	if locationParams.Latitude < -90 || locationParams.Latitude > 90 || locationParams.Latitude == 0 {
		return storage.ErrInvalidLatitude
	}
	if locationParams.Longitude < -180 || locationParams.Longitude > 180 || locationParams.Longitude == 0 {
		return storage.ErrInvalidLongitude
	}

	return nil
}
