package tree

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/utils"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
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
	log := logger.GetLogger(ctx)
	log.Debug("handle event", "event", event.Type(), "service", "TreeService")
	t, err := s.treeRepo.GetBySensorID(ctx, event.New.SensorID)
	if err != nil {
		log.Error("failed to get tree by sensor id", "sensor_id", event.New.SensorID, "err", err)
		return nil
	}

	status := utils.CalculateWateringStatus(t.PlantingYear, event.New.Data.Watermarks)

	if status == t.WateringStatus {
		return nil
	}

	newTree, err := s.treeRepo.Update(ctx, t.ID, tree.WithWateringStatus(status))
	if err != nil {
		log.Error("failed to update tree with new watering status", "tree_id", t.ID, "watering_status", status, "err", err)
	}

	slog.Info("updating tree watering status", "prev_status", t.WateringStatus, "new_status", status)

	s.publishUpdateTreeEvent(ctx, t, newTree)
	return nil
}

func (s *TreeService) GetAll(ctx context.Context) ([]*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	trees, err := s.treeRepo.GetAll(ctx)
	if err != nil {
		log.Error("failed to fetch trees", "error", err)
		return nil, handleError(err)
	}

	return trees, nil
}

func (s *TreeService) GetByID(ctx context.Context, id int32) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	tr, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to fetch tree by id", "error", err, "tree_id", id)
		return nil, handleError(err)
	}

	return tr, nil
}

func (s *TreeService) GetBySensorID(ctx context.Context, id string) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	tr, err := s.treeRepo.GetBySensorID(ctx, id)
	if err != nil {
		log.Error("failed to get tree by sensor id", "sensor_id", id, "error", err)
		return nil, handleError(err)
	}

	return tr, nil
}

func (s *TreeService) publishUpdateTreeEvent(ctx context.Context, prevTree, updatedTree *entities.Tree) {
	log := logger.GetLogger(ctx)
	log.Debug("publish new event", "event", entities.EventTypeUpdateTree, "service", "TreeService")
	event := entities.NewEventUpdateTree(prevTree, updatedTree)
	if err := s.eventManager.Publish(ctx, event); err != nil {
		log.Error("error while sending event after updating tree", "err", err, "tree_id", prevTree.ID)
	}
}

func (s *TreeService) publishCreateTreeEvent(ctx context.Context, newTree *entities.Tree) {
	log := logger.GetLogger(ctx)
	log.Debug("publish new event", "event", entities.EventTypeCreateTree, "service", "TreeService")
	event := entities.NewEventCreateTree(newTree)
	if err := s.eventManager.Publish(ctx, event); err != nil {
		log.Error("error while sending event after creating tree", "err", err, "tree_id", newTree.ID)
	}
}

func (s *TreeService) publishDeleteTreeEvent(ctx context.Context, prevTree *entities.Tree) {
	log := logger.GetLogger(ctx)
	log.Debug("publish new event", "event", entities.EventTypeDeleteTree, "service", "TreeService")
	event := entities.NewEventDeleteTree(prevTree)
	if err := s.eventManager.Publish(ctx, event); err != nil {
		log.Error("error while sending event after deleting tree", "err", err, "tree_id", prevTree.ID)
	}
}

func (s *TreeService) Create(ctx context.Context, treeCreate *entities.TreeCreate) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	if err := s.validator.Struct(treeCreate); err != nil {
		log.Debug("failed to validate tree struct to create", "error", err, "raw_tree", fmt.Sprintf("%+v", treeCreate))
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	fn := make([]entities.EntityFunc[entities.Tree], 0)
	if treeCreate.TreeClusterID != nil {
		treeClusterID, err := s.treeClusterRepo.GetByID(ctx, *treeCreate.TreeClusterID)
		if err != nil {
			log.Error("failed to fetch tree cluster by id specified in the tree create request", "tree_cluster_id", treeCreate.TreeClusterID)
			return nil, handleError(err)
		}
		fn = append(fn, tree.WithTreeCluster(treeClusterID))
	}

	if treeCreate.SensorID != nil {
		sensorID, err := s.sensorRepo.GetByID(ctx, *treeCreate.SensorID)
		if err != nil {
			log.Error("failed to fetch sensor by id specified in the tree create request", "sensor_id", treeCreate.SensorID)
			return nil, handleError(err)
		}
		fn = append(fn, tree.WithSensor(sensorID))
	}

	fn = append(fn,
		tree.WithReadonly(treeCreate.Readonly),
		tree.WithPlantingYear(treeCreate.PlantingYear),
		tree.WithSpecies(treeCreate.Species),
		tree.WithNumber(treeCreate.Number),
		tree.WithLatitude(treeCreate.Latitude),
		tree.WithLongitude(treeCreate.Longitude),
	)
	newTree, err := s.treeRepo.Create(ctx, fn...)
	if err != nil {
		log.Error("failed to create tree", "error", err)
		return nil, handleError(err)
	}

	slog.Info("successfully created tree", "tree_id", newTree.ID)
	s.publishCreateTreeEvent(ctx, newTree)
	return newTree, nil
}

func (s *TreeService) Delete(ctx context.Context, id int32) error {
	log := logger.GetLogger(ctx)
	treeEntity, err := s.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}
	if err := s.treeRepo.Delete(ctx, id); err != nil {
		log.Error("failed to delete tree", "error", err, "tree_id", id)
		return handleError(err)
	}

	slog.Info("successfully deleting tree", "tree_id", id)
	s.publishDeleteTreeEvent(ctx, treeEntity)
	return nil
}

func (s *TreeService) Update(ctx context.Context, id int32, tu *entities.TreeUpdate) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	if err := s.validator.Struct(tu); err != nil {
		log.Debug("failed to validate struct from tree update", "error", err, "raw_tree", fmt.Sprintf("%+v", tu))
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	prevTree, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	// Check if the tree is readonly (imported from csv)
	// if currentTree.Readonly {
	// 	return nil, handleError(fmt.Errorf("tree with ID %d is readonly and cannot be updated", id))
	// }

	fn := make([]entities.EntityFunc[entities.Tree], 0)
	if tu.TreeClusterID != nil {
		var treeCluster *entities.TreeCluster
		treeCluster, err = s.treeClusterRepo.GetByID(ctx, *tu.TreeClusterID)
		if err != nil {
			log.Error("failed to find tree cluster by id specified from update request", "tree_cluster_id", tu.TreeClusterID)
			return nil, handleError(fmt.Errorf("failed to find TreeCluster with ID %d: %w", *tu.TreeClusterID, err))
		}
		fn = append(fn, tree.WithTreeCluster(treeCluster))
	} else {
		fn = append(fn, tree.WithTreeCluster(nil))
	}

	if tu.SensorID != nil {
		var sensor *entities.Sensor
		sensor, err = s.sensorRepo.GetByID(ctx, *tu.SensorID)
		if err != nil {
			log.Error("failed to find sensor by id specified from update request", "sensor_id", tu.SensorID)
			return nil, handleError(fmt.Errorf("failed to find Sensor with ID %v: %w", *tu.SensorID, err))
		}
		fn = append(fn, tree.WithSensor(sensor))
	} else {
		fn = append(fn, tree.WithSensor(nil))
	}

	fn = append(fn, tree.WithPlantingYear(tu.PlantingYear),
		tree.WithSpecies(tu.Species),
		tree.WithNumber(tu.Number),
		tree.WithLatitude(tu.Latitude),
		tree.WithLongitude(tu.Longitude),
		tree.WithDescription(tu.Description))

	updatedTree, err := s.treeRepo.Update(ctx, id, fn...)
	if err != nil {
		log.Error("failed to update tree", "error", err, "tree_id", id)
		return nil, handleError(err)
	}

	slog.Info("successfully updated tree", "tree_id", id)
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
