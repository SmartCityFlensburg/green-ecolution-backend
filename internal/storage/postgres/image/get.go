package image

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (r *ImageRepository) GetAll(ctx context.Context) ([]*entities.Image, error) {
	rows, err := r.store.GetAllImages(ctx)
	if err != nil {
		slog.Error("Error getting all images", "Error", err)
		return nil, err
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *ImageRepository) GetByID(ctx context.Context, id int32) (*entities.Image, error) {
	row, err := r.store.GetImageByID(ctx, id)
	if err != nil {
		slog.Error("Error getting image by ID", "Error", err, "Image ID", id)
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}
