package region

import (
	"context"
	"errors"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (r *RegionRepository) Update(ctx context.Context, id int32, vFn ...entities.EntityFunc[entities.Region]) (*entities.Region, error) {
	entity, err := r.GetByID(ctx, id)
	if err != nil {
		slog.Error("Error getting region by id", "Error", err, "ID", id)
		return nil, err
	}

	for _, fn := range vFn {
		fn(entity)
	}

	if entity.Name == "" {
		return nil, errors.New("name is required")
	}

	if err := r.updateEntity(ctx, entity); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, entity.ID)
}

func (r *RegionRepository) updateEntity(ctx context.Context, vehicle *entities.Region) error {
	params := sqlc.UpdateRegionParams{
		ID:   vehicle.ID,
		Name: vehicle.Name,
	}

	return r.store.UpdateRegion(ctx, &params)
}
