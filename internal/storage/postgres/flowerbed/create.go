package flowerbed

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/pkg/errors"
)

func defaultFlowerbed() *entities.Flowerbed {
	return &entities.Flowerbed{
		Sensor:         &entities.Sensor{},
		Size:           0,
		Description:    "",
		NumberOfPlants: 0,
		MoistureLevel:  0,
		Region:         nil,
		Address:        "",
		Latitude:       nil,
		Longitude:      nil,
		Images:         make([]*entities.Image, 0),
		Archived:       false,
	}
}

func (r *FlowerbedRepository) Create(ctx context.Context, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error) {
	entity := defaultFlowerbed()
	for _, fn := range fFn {
		fn(entity)
	}

	if err := r.validateFlowerbedEntity(entity); err != nil {
		return nil, err
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		slog.Error("Error creating flowerbed", "Error", err)
		return nil, r.store.HandleError(err)
	}

	entity.ID = *id

	return r.GetByID(ctx, *id)
}

func (r *FlowerbedRepository) CreateAndLinkImages(ctx context.Context, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error) {
	entity, err := r.Create(ctx, fFn...)
	if err != nil {
		slog.Error("Error creating ", "Error", err)
		return nil, err
	}

	if err := r.handleImages(ctx, entity.ID, entity.Images); err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *FlowerbedRepository) createEntity(ctx context.Context, entity *entities.Flowerbed) (*int32, error) {
	var region *int32
	if entity.Region != nil {
		region = &entity.Region.ID
	}

	args := sqlc.CreateFlowerbedParams{
		SensorID:       &entity.Sensor.ID,
		Size:           entity.Size,
		Description:    entity.Description,
		NumberOfPlants: entity.NumberOfPlants,
		MoistureLevel:  entity.MoistureLevel,
		RegionID:       region,
		Address:        entity.Address,
		Latitude:       *entity.Latitude,
		Longitude:      *entity.Longitude,
	}

	id, err := r.store.CreateFlowerbed(ctx, &args)
	if err != nil {
		slog.Error("Error creating flowerbed", "Error", err)
		return nil, err
	}

	return &id, nil
}

func (r *FlowerbedRepository) handleImages(ctx context.Context, flowerbedID int32, images []*entities.Image) error {
	for _, img := range images {
		err := r.linkImages(ctx, flowerbedID, img.ID)
		if err != nil {
			slog.Error("Error linking images to flowerbed", "Error", err)
			return err
		}
	}

	return nil
}

func (r *FlowerbedRepository) linkImages(ctx context.Context, flowerbedID, imgID int32) error {
	params := sqlc.LinkFlowerbedImageParams{
		FlowerbedID: flowerbedID,
		ImageID:     imgID,
	}
	return r.store.LinkFlowerbedImage(ctx, &params)
}

func (r *FlowerbedRepository) validateFlowerbedEntity(fb *entities.Flowerbed) error {
	if fb == nil {
		return errors.New("flowerbed is nil")
	}

	if fb.Longitude == nil {
		return errors.New("flowerbed longitude is empty")
	}

	if fb.Latitude == nil {
		return errors.New("flowerbed latitude is empty")
	}

	if fb.Region == nil {
		return errors.New("flowerbed region is empty")
	}

	return nil
}
