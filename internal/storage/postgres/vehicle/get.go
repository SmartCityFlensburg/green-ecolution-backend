package vehicle

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (r *VehicleRepository) GetAll(ctx context.Context) ([]*entities.Vehicle, error) {
	rows, err := r.store.GetAllVehicles(ctx)
	if err != nil {
		slog.Error("Error getting all vehicles", "Error", err)
		return nil, err
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *VehicleRepository) GetAllByType(ctx context.Context, vehicleType entities.VehicleType) ([]*entities.Vehicle, error) {
	rows, err := r.store.GetAllVehiclesByType(ctx, sqlc.VehicleType(vehicleType))
	if err != nil {
		slog.Error("Error getting all vehicles by type", "Error", err, "VehicleType", vehicleType)
		return nil, err
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *VehicleRepository) GetByID(ctx context.Context, id int32) (*entities.Vehicle, error) {
	row, err := r.store.GetVehicleByID(ctx, id)
	if err != nil {
		slog.Error("Error getting vehicle by ID", "Error", err, "VehicleID", id)
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *VehicleRepository) GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error) {
	row, err := r.store.GetVehicleByPlate(ctx, plate)
	if err != nil {
		slog.Error("Error getting vehicle by plate", "Error", err, "Plate", plate)
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}
