package tree

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/utils"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker"
)

type TreeService struct {
	treeRepo        storage.TreeRepository
	sensorRepo      storage.SensorRepository
	ImageRepo       storage.ImageRepository
	treeClusterRepo storage.TreeClusterRepository
	validator       *validator.Validate
	eventManager    *worker.EventManager
}

func NewTreeService(
	repoTree storage.TreeRepository,
	repoSensor storage.SensorRepository,
	repoImage storage.ImageRepository,
	treeClusterRepo storage.TreeClusterRepository,
	eventManager *worker.EventManager,
) service.TreeService {
	return &TreeService{
		treeRepo:        repoTree,
		sensorRepo:      repoSensor,
		ImageRepo:       repoImage,
		treeClusterRepo: treeClusterRepo,
		validator:       validator.New(),
		eventManager:    eventManager,
	}
}

func (s *TreeService) HandleNewSensorData(ctx context.Context, event *entities.EventNewSensorData) error {
	slog.Debug("handle event", "event", event.Type(), "service", "TreeService")
	t, err := s.treeRepo.GetBySensorID(ctx, event.New.SensorID)
	if err != nil {
		slog.Error("failed to get tree by sensor id", "sensor_id", event.New.SensorID, "err", err)
		return nil
	}

	status := utils.CalculateWateringStatus(t.PlantingYear, event.New.Data.Watermarks)

	if status == t.WateringStatus {
		return nil
	}

	newTree, err := s.treeRepo.Update(ctx, t.ID, func(tree *entities.Tree) (bool, error) {
		tree.WateringStatus = status
		return true, nil
	})

	if err != nil {
		slog.Error("failed to update tree with new watering status", "tree_id", t.ID, "watering_status", status, "err", err)
	}

	s.publishUpdateTreeEvent(ctx, t, newTree)
	return nil
}

func (s *TreeService) GetAll(ctx context.Context) ([]*entities.Tree, error) {
	trees, err := s.treeRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return trees, nil
}

func (s *TreeService) GetByID(ctx context.Context, id int32) (*entities.Tree, error) {
	tr, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return tr, nil
}

func (s *TreeService) GetBySensorID(ctx context.Context, id string) (*entities.Tree, error) {
	tr, err := s.treeRepo.GetBySensorID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return tr, nil
}

func (s *TreeService) publishUpdateTreeEvent(ctx context.Context, prevTree, updatedTree *entities.Tree) {
	slog.Debug("publish new event", "event", entities.EventTypeUpdateTree, "service", "TreeService")
	event := entities.NewEventUpdateTree(prevTree, updatedTree)
	if err := s.eventManager.Publish(ctx, event); err != nil {
		slog.Error("error while sending event after updating tree", "err", err, "tree_id", prevTree.ID)
	}
}

func (s *TreeService) publishCreateTreeEvent(ctx context.Context, newTree *entities.Tree) {
	slog.Debug("publish new event", "event", entities.EventTypeCreateTree, "service", "TreeService")
	event := entities.NewEventCreateTree(newTree)
	if err := s.eventManager.Publish(ctx, event); err != nil {
		slog.Error("error while sending event after creating tree", "err", err, "tree_id", newTree.ID)
	}
}

func (s *TreeService) publishDeleteTreeEvent(ctx context.Context, prevTree *entities.Tree) {
	slog.Debug("publish new event", "event", entities.EventTypeDeleteTree, "service", "TreeService")
	event := entities.NewEventDeleteTree(prevTree)
	if err := s.eventManager.Publish(ctx, event); err != nil {
		slog.Error("error while sending event after deleting tree", "err", err, "tree_id", prevTree.ID)
	}
}

func (s *TreeService) Create(ctx context.Context, treeCreate *entities.TreeCreate) (*entities.Tree, error) {
	if err := s.validator.Struct(treeCreate); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	var treeClusterID *entities.TreeCluster
	if treeCreate.TreeClusterID != nil {
		var err error
		treeClusterID, err = s.treeClusterRepo.GetByID(ctx, *treeCreate.TreeClusterID)
		if err != nil {
			return nil, handleError(err)
		}
	}

	var sensorID *entities.Sensor
	if treeCreate.SensorID != nil {
		var err error
		sensorID, err = s.sensorRepo.GetByID(ctx, *treeCreate.SensorID)
		if err != nil {
			return nil, handleError(err)
		}
	}

	newTree, err := s.treeRepo.Create(ctx, func(tree *entities.Tree) (bool, error) {
		tree.Readonly = treeCreate.Readonly
		tree.PlantingYear = treeCreate.PlantingYear
		tree.Species = treeCreate.Species
		tree.Number = treeCreate.Number
		tree.Latitude = treeCreate.Latitude
		tree.Longitude = treeCreate.Longitude

		// Apply TreeCluster and Sensor if available
		if treeClusterID != nil {
			tree.TreeCluster = treeClusterID
		}
		if sensorID != nil {
			tree.Sensor = sensorID
		}

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	s.publishCreateTreeEvent(ctx, newTree)
	return newTree, nil
}

func (s *TreeService) Delete(ctx context.Context, id int32) error {
	treeEntity, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}
	if err := s.treeRepo.Delete(ctx, id); err != nil {
		return handleError(err)
	}

	s.publishDeleteTreeEvent(ctx, treeEntity)
	return nil
}

func (s *TreeService) Update(ctx context.Context, id int32, tu *entities.TreeUpdate) (*entities.Tree, error) {
	if err := s.validator.Struct(tu); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	prevTree, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	var treeCluster *entities.TreeCluster
	if tu.TreeClusterID != nil {
		treeCluster, err = s.treeClusterRepo.GetByID(ctx, *tu.TreeClusterID)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to find TreeCluster with ID %d: %w", *tu.TreeClusterID, err))
		}
	}

	var sensor *entities.Sensor
	if tu.SensorID != nil {
		sensor, err = s.sensorRepo.GetByID(ctx, *tu.SensorID)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to find Sensor with ID %v: %w", *tu.SensorID, err))
		}
	}

	updatedTree, err := s.treeRepo.Update(ctx, id, func(tree *entities.Tree) (bool, error) {
		tree.PlantingYear = tu.PlantingYear
		tree.Species = tu.Species
		tree.Number = tu.Number
		tree.Latitude = tu.Latitude
		tree.Longitude = tu.Longitude
		tree.Description = tu.Description
		tree.TreeCluster = treeCluster
		tree.Sensor = sensor
		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	s.publishUpdateTreeEvent(ctx, prevTree, updatedTree)
	return updatedTree, nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, storage.ErrTreeNotFound.Error())
	}

	if errors.Is(err, storage.ErrSensorNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (s *TreeService) Ready() bool {
	return s.treeRepo != nil && s.sensorRepo != nil
}
