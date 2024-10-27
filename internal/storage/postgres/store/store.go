package store

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type EntityType string

const (
	Sensor      EntityType = "sensor"
	Image       EntityType = "image"
	Flowerbed   EntityType = "flowerbed"
	TreeCluster EntityType = "treecluster"
	Tree        EntityType = "tree"
	Vehicle     EntityType = "vehicle"
)

type Store struct {
	*sqlc.Queries
	db         *pgxpool.Pool
	entityType EntityType
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: sqlc.New(db),
		db:      db,
	}
}

func (s *Store) SetEntityType(entityType EntityType) {
	s.entityType = entityType
}

func (s *Store) HandleError(err error) error {
	if err == nil {
		return nil
	}

	slog.Error("An Error occurred in database operation", "error", err)
	return err
}

func (s *Store) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		return err
	}

	err = fn(tx)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			slog.Error("Failed to rollback transaction", "rollbackError", rbErr, "originalError", err)
			return rbErr
		}
		slog.Error("Transaction function failed", "error", err)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		return err

	}
	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) CheckSensorExists(ctx context.Context, sensorID *int32) error {
	if sensorID != nil {
		_, err := s.GetSensorByID(ctx, *sensorID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return storage.ErrSensorNotFound
			} else {
				slog.Error("Error getting sensor by id", "sensorID", *sensorID, "error", err)
				return s.HandleError(err)
			}
		}
	}

	return nil
}

func (s *Store) CheckImageExists(ctx context.Context, imageID *int32) error {
	if imageID != nil {
		_, err := s.GetImageByID(ctx, *imageID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return storage.ErrImageNotFound
			} else {
				slog.Error("Error getting image by id", "imageID", *imageID, "error", err)
				return s.HandleError(err)
			}
		}
	}

	return nil
}

func (s *Store) CheckTreeClusterExists(ctx context.Context, treeClusterID *int32) error {
	if treeClusterID != nil {
		_, err := s.GetTreeClusterByID(ctx, *treeClusterID)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "42P01" {
				slog.Error("Database table 'tree_clusters' does not exist", "error", err)
				return errors.New("database table 'tree_clusters' does not exist")
			}
			if errors.Is(err, pgx.ErrNoRows) {
				return storage.ErrTreeClusterNotFound
			}
			slog.Error("Error getting tree cluster by id", "treeClusterID", *treeClusterID, "error", err)
			return s.HandleError(err)

		}

	}
	return nil
}
