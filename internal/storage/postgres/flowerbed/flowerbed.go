package flowerbed

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/pkg/errors"
)

type FlowerbedRepository struct {
	store *store.Store
	FlowerbedMappers
}

type FlowerbedMappers struct {
	mapper       mapper.InternalFlowerbedRepoMapper
	imgMapper    mapper.InternalImageRepoMapper
	sensorMapper mapper.InternalSensorRepoMapper
	regionMapper mapper.InternalRegionRepoMapper
}

func NewFlowerbedMappers(
	fMapper mapper.InternalFlowerbedRepoMapper,
	iMapper mapper.InternalImageRepoMapper,
	sMapper mapper.InternalSensorRepoMapper,
	rMapper mapper.InternalRegionRepoMapper,
) FlowerbedMappers {
	return FlowerbedMappers{
		mapper:       fMapper,
		imgMapper:    iMapper,
		sensorMapper: sMapper,
		regionMapper: rMapper,
	}
}

func NewFlowerbedRepository(s *store.Store, mappers FlowerbedMappers) storage.FlowerbedRepository {
	return &FlowerbedRepository{
		store:            s,
		FlowerbedMappers: mappers,
	}
}

func WithSize(size float64) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Size = size
	}
}

func WithDescription(description string) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Description = description
	}
}

func WithNumberOfPlants(numberOfPlants int32) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.NumberOfPlants = numberOfPlants
	}
}

func WithMoistureLevel(moistureLevel float64) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.MoistureLevel = moistureLevel
	}
}

func WithRegion(region *entities.Region) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Region = region
	}
}

func WithAddress(address string) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Address = address
	}
}

func WithLatitude(latitude *float64) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Latitude = latitude
	}
}

func WithLongitude(longitude *float64) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Longitude = longitude
	}
}

func WithSensor(sensor *entities.Sensor) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Sensor = sensor
	}
}

func WithImages(images []*entities.Image) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Images = images
	}
}

func WithArchived(archived bool) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Archived = archived
	}
}

func WithImagesIDs(imagesIDs []int32) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		for _, id := range imagesIDs {
			f.Images = append(f.Images, &entities.Image{ID: id})
		}
	}
}

func WithSensorID(id string) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Sensor = &entities.Sensor{ID: id}
	}
}

func WithRegionID(id int32) entities.EntityFunc[entities.Flowerbed] {
	return func(f *entities.Flowerbed) {
		f.Region = &entities.Region{ID: id}
	}
}

func (r *FlowerbedRepository) Delete(ctx context.Context, id int32) error {
	log := logger.GetLogger(ctx)
	_, err := r.store.DeleteFlowerbed(ctx, id)
	if err != nil {
		log.Error("failed to delete flowerbed entity in db", "error", err, "flowerbed_id", id)
		return err
	}

	log.Debug("flowerbed entity deleted successfully in db", "flowerbed_id", id)
	return nil
}

func (r *FlowerbedRepository) DeleteAndUnlinkImages(ctx context.Context, id int32) error {
	images, err := r.GetAllImagesByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get images")
	}

	for _, img := range images {
		if err := r.UnlinkImage(ctx, id, img.ID); err != nil {
			return errors.Wrap(err, "failed to unlink images")
		}
	}

	return r.Delete(ctx, id)
}

func (r *FlowerbedRepository) UnlinkImage(ctx context.Context, id, imageID int32) error {
	args := sqlc.UnlinkFlowerbedImageParams{
		FlowerbedID: id,
		ImageID:     imageID,
	}

	_, err := r.store.UnlinkFlowerbedImage(ctx, &args)
	if err != nil {
		return err
	}

	return nil
}

func (r *FlowerbedRepository) UnlinkAllImages(ctx context.Context, id int32) error {
	_, err := r.store.UnlinkAllFlowerbedImages(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *FlowerbedRepository) UnlinkSensorID(ctx context.Context, sensorID string) error {
	if sensorID == "" {
		return errors.New("sensorID cannot be empty")
	}
	return r.store.UnlinkSensorIDFromFlowerbeds(ctx, &sensorID)
}

func (r *FlowerbedRepository) Archive(ctx context.Context, id int32) error {
	log := logger.GetLogger(ctx)
	_, err := r.store.ArchiveFlowerbed(ctx, id)
	if err != nil {
		log.Error("failed to archive flowerbed entity in db", "error", err, "flowerbed_id", id)
		return err
	}

	log.Debug("flowerbed entity archived successfully in db", "flowerbed_id", id)
	return nil
}
