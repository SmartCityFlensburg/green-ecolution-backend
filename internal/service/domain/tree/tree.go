package tree

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
)

type TreeService struct {
	treeRepo        storage.TreeRepository
	sensorRepo      storage.SensorRepository
	treeClusterRepo storage.TreeClusterRepository
	validator       *validator.Validate
}

func NewTreeService(repoTree storage.TreeRepository, repoSensor storage.SensorRepository, treeClusterRepo storage.TreeClusterRepository) service.TreeService {
	return &TreeService{
		treeRepo:        repoTree,
		sensorRepo:      repoSensor,
		treeClusterRepo: treeClusterRepo,
		validator:       validator.New(),
	}
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

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (s *TreeService) Ready() bool {
	return s.treeRepo != nil && s.sensorRepo != nil
}
func (s *TreeService) Create(ctx context.Context, treeCreate *entities.TreeCreate) (*entities.Tree, error) {
	if treeCreate.PlantingYear == 0 {
		return nil, handleError(errors.New("plantingYear cannot be null or zero"))
	}
	if treeCreate.Species == "" {
		return nil, handleError(errors.New("species (Gattung) cannot be null or empty"))
	}
	if treeCreate.Number == "" {
		return nil, handleError(errors.New("tree Number (Baum Nr) cannot be null or empty"))
	}
	if treeCreate.Latitude == 0 || treeCreate.Longitude == 0 {
		return nil, handleError(errors.New("latitude and Longitude cannot be null or zero"))
	}

	fn := make([]entities.EntityFunc[entities.Tree], 0)
	if treeCreate.TreeClusterID != nil {
		treeClusterID, err := s.treeClusterRepo.GetByID(ctx, *treeCreate.TreeClusterID)
		if err != nil {
			return nil, handleError(err)
		}
		fn = append(fn, tree.WithTreeCluster(treeClusterID))
	}
	fn = append(fn,
		tree.WithReadonly(treeCreate.Readonly),
		tree.WithPlantingYear(treeCreate.PlantingYear),
		tree.WithSpecies(treeCreate.Species),
		tree.WithTreeNumber(treeCreate.Number),
		tree.WithLatitude(treeCreate.Latitude),
		tree.WithLongitude(treeCreate.Longitude),
	)
	newTree, err := s.treeRepo.Create(ctx, fn...)
	if err != nil {
		return nil, handleError(err)
	}
	// TODO: update the coordinates of the tree cluster.
	return newTree, nil
}

func (s *TreeService) Delete(ctx context.Context, id int32) error {
	_, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}
	err = s.treeRepo.DeleteAndUnlinkImages(ctx, id)
	if err != nil {
		return handleError(err)
	}
	// TODO: update the coordinates of the tree cluster.
	return nil
}

func (s *TreeService) Update(ctx context.Context, id int32, tu *entities.TreeUpdate) (*entities.Tree, error) {
	err := s.validator.RegisterValidation("not-zero", notZero)
	if err != nil {
		return nil, handleError(err)
	}
	if err := s.validator.Struct(tu); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}
	currentTree, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}
	// Check if the tree is readonly (imported from csv)
	if currentTree.Readonly {
		return nil, handleError(fmt.Errorf("tree with ID %d is readonly and cannot be updated", id))
	}
	fn := make([]entities.EntityFunc[entities.Tree], 0)

	if tu.TreeClusterID != nil {
		treeCluster, err := s.treeClusterRepo.GetByID(ctx, *tu.TreeClusterID)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to find TreeCluster with ID %d: %w", *tu.TreeClusterID, err))
		}
		fn = append(fn, tree.WithTreeCluster(treeCluster))
	}
	fn = append(fn, tree.WithPlantingYear(tu.PlantingYear),
		tree.WithSpecies(tu.Species),
		tree.WithTreeNumber(tu.Number),
		tree.WithLatitude(tu.Latitude),
		tree.WithLongitude(tu.Longitude))
	updatedTree, err := s.treeRepo.Update(ctx, id, fn...)
	if err != nil {
		return nil, handleError(err)
	}
	// TODO: If a new tree cluster has been provided, update the coordinates of both the old tree cluster and the new one.
	return updatedTree, nil
}
func notZero(fl validator.FieldLevel) bool {
	value := fl.Field().Float()
	return value != 0
}
