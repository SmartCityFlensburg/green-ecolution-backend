package tree

import (
	"context"
	"log/slog"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (r *TreeRepository) GetAll(ctx context.Context) ([]*entities.Tree, error) {
	rows, err := r.store.GetAllTrees(ctx)
	if err != nil {
		slog.Error("Error getting all trees", "Error", err)
		return nil, err
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			slog.Error("Error mapping fields", "Error", err)
			return nil, err
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByID(ctx context.Context, id int32) (*entities.Tree, error) {
	row, err := r.store.GetTreeByID(ctx, id)
	if err != nil {
		slog.Error("Error getting tree by id", "Error", err)
		return nil, err
	}

	t := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, t); err != nil {
		return nil, r.store.HandleError(err)
	}

	return t, nil
}

func (r *TreeRepository) GetBySensorID(ctx context.Context, id string) (*entities.Tree, error) {
	_, err := r.store.GetSensorByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrSensorNotFound
		} else {
			slog.Error("Error getting sensor by id", "Error", err, "SensorID", id)
			return nil, err
		}
	}

	row, err := r.store.GetTreeBySensorID(ctx, &id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, t); err != nil {
		return nil, r.store.HandleError(err)
	}

	return t, nil
}

func (r *TreeRepository) GetBySensorIDs(ctx context.Context, ids ...string) ([]*entities.Tree, error) {
	rows, err := r.store.GetTreesBySensorIDs(ctx, ids)
	if err != nil {
		slog.Error("Error getting trees by sensor ids", "Error", err, "SensorIDs", ids)
		return nil, err
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return t, nil
}

func (r *TreeRepository) GetTreesByIDs(ctx context.Context, ids []int32) ([]*entities.Tree, error) {
	rows, err := r.store.GetTreesByIDs(ctx, ids)
	if err != nil {
		slog.Error("Error getting trees by ids", "Error", err, "TreeIDs", ids)
		return nil, err
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error) {
	_, err := r.store.GetTreeClusterByID(ctx, id)
	if err != nil {
		slog.Error("Error getting tree cluster by id", "Error", err, "TreeClusterID", id)
		return nil, err
	}

	rows, err := r.store.GetTreesByTreeClusterID(ctx, &id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByCoordinates(ctx context.Context, latitude, longitude float64) (*entities.Tree, error) {
	params := sqlc.GetTreeByCoordinatesParams{
		Latitude:  latitude,
		Longitude: longitude,
	}
	row, err := r.store.GetTreeByCoordinates(ctx, &params)
	if err != nil {
		slog.Error("Error getting tree by coordinates", "Error", err, "Latitude", latitude, "Longitude", longitude)
		return nil, err
	}
	tree := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, tree); err != nil {
		return nil, r.store.HandleError(err)
	}
	return tree, nil
}

func (r *TreeRepository) GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error) {
	rows, err := r.store.GetAllImagesByTreeID(ctx, id)
	if err != nil {
		slog.Error("Error getting all images by tree id", "Error", err, "TreeID", id)
		return nil, err
	}

	return r.iMapper.FromSqlList(rows), nil
}

func (r *TreeRepository) GetSensorByTreeID(ctx context.Context, treeID int32) (*entities.Sensor, error) {
	row, err := r.store.GetSensorByTreeID(ctx, treeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrSensorNotFound
		} else {
			slog.Error("Error getting sensor by tree id", "Error", err, "TreeID", treeID)
			return nil, err
		}
	}

	data := r.sMapper.FromSql(row)
	if err := r.store.MapSensorFields(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *TreeRepository) getTreeClusterByTreeID(ctx context.Context, treeID int32) (*entities.TreeCluster, error) {
	row, err := r.store.GetTreeClusterByTreeID(ctx, treeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrTreeClusterNotFound
		} else {
			slog.Error("Error getting tree cluster by tree id", "Error", err, "TreeID", treeID)
			return nil, err
		}
	}

	return r.tcMapper.FromSql(row), nil
}

// Map sensor, images and tree cluster entity to domain flowerbed
func (r *TreeRepository) mapFields(ctx context.Context, t *entities.Tree) error {
	if err := mapImages(ctx, r, t); err != nil {
		slog.Error("Error mapping images", "Error", err)
		return err
	}

	if err := mapSensor(ctx, r, t); err != nil {
		return r.store.HandleError(err)
	}

	_ = mapTreeCluster(ctx, r, t)

	return nil
}

func mapImages(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	images, err := r.GetAllImagesByID(ctx, t.ID)
	if err != nil {
		slog.Error("Error getting all images by tree id", "Error", err, "TreeID", t.ID)
		return err
	}
	t.Images = images
	return nil
}

func mapSensor(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	sensor, err := r.GetSensorByTreeID(ctx, t.ID)
	if err != nil {
		if errors.Is(err, storage.ErrSensorNotFound) {
			// If sensor is not found, set sensor to nil
			t.Sensor = nil
			return nil
		}
		slog.Error("Error getting sensor by tree id", "Error", err, "TreeID", t.ID)
		return err
	}
	t.Sensor = sensor
	return nil
}

func mapTreeCluster(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	treeCluster, err := r.getTreeClusterByTreeID(ctx, t.ID)
	if err != nil {
		slog.Error("Error getting tree cluster by tree id", "Error", err, "TreeID", t.ID)
		return err
	}
	t.TreeCluster = treeCluster
	return nil
}

func (r *TreeRepository) FindNearestTree(ctx context.Context, latitude, longitude float64) (*entities.Tree, error) {
	params := &sqlc.FindNearestTreeParams{
		StMakepoint:   latitude,
		StMakepoint_2: longitude,
	}

	nearestTree, err := r.store.FindNearestTree(ctx, params)
	if err != nil {
		return nil, err
	}

	tree := r.mapper.FromSql(nearestTree)
	if err := r.mapFields(ctx, tree); err != nil {
		slog.Error("Error mapping fields", "Error", err)
		return nil, err
	}
	return tree, nil
}
