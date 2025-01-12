package wateringplan

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func defaultWateringPlan() *entities.WateringPlan {
	return &entities.WateringPlan{
		Date:               time.Time{},
		Description:        "",
		Distance:           utils.P(0.0),
		TotalWaterRequired: utils.P(0.0),
		Status:             entities.WateringPlanStatusPlanned,
		UserIDs:            make([]*uuid.UUID, 0),
		TreeClusters:       make([]*entities.TreeCluster, 0),
		Transporter:        nil,
		Trailer:            nil,
		CancellationNote:   "",
		Duration:           time.Duration(0),
		RefillCount:        0,
	}
}

func (w *WateringPlanRepository) Create(ctx context.Context, createFn func(*entities.WateringPlan) (bool, error)) (*entities.WateringPlan, error) {
	if createFn == nil {
		return nil, errors.New("createFn is nil")
	}

	var createdWp *entities.WateringPlan
	err := w.store.WithTx(ctx, func(s *store.Store) error {
		oldStore := w.store
		defer func() {
			w.store = oldStore
		}()
		w.store = s

		entity := defaultWateringPlan()
		created, err := createFn(entity)
		if err != nil {
			return err
		}

		if !created {
			return nil
		}

		if err := w.validateWateringPlan(entity); err != nil {
			return err
		}

		id, err := w.createEntity(ctx, entity)
		if err != nil {
			return err
		}

		createdWp, err = w.GetByID(ctx, *id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdWp, nil
}

func (w *WateringPlanRepository) createEntity(ctx context.Context, entity *entities.WateringPlan) (*int32, error) {
	date, err := utils.TimeToPgDate(entity.Date)
	if err != nil {
		return nil, errors.New("failed to convert date")
	}

	args := sqlc.CreateWateringPlanParams{
		Date:               date,
		Description:        entity.Description,
		Distance:           entity.Distance,
		TotalWaterRequired: entity.TotalWaterRequired,
		Status:             sqlc.WateringPlanStatus(entities.WateringPlanStatusPlanned),
	}

	id, err := w.store.CreateWateringPlan(ctx, &args)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	if err := w.setLinkedUsers(ctx, entity, id); err != nil {
		return nil, w.store.HandleError(err)
	}

	if err := w.setLinkedVehicles(ctx, entity, id); err != nil {
		return nil, w.store.HandleError(err)
	}

	if err := w.setLinkedTreeClusters(ctx, entity, id); err != nil {
		return nil, w.store.HandleError(err)
	}

	return &id, nil
}

func (w *WateringPlanRepository) validateWateringPlan(entity *entities.WateringPlan) error {
	if entity.Transporter == nil || entity.Transporter.Type != entities.VehicleTypeTransporter {
		return errors.New("watering plan requires a valid transporter")
	}

	if entity.Trailer != nil && entity.Trailer.Type != entities.VehicleTypeTrailer {
		return errors.New("trailer vehicle requires a vehicle of type trailer")
	}

	if len(entity.UserIDs) == 0 {
		return errors.New("watering plan requires employees")
	}

	if len(entity.TreeClusters) == 0 {
		return errors.New("watering plan requires tree cluster")
	}

	return nil
}

func (w *WateringPlanRepository) setLinkedVehicles(ctx context.Context, entity *entities.WateringPlan, id int32) error {
	// link transporter to pivot table
	err := w.store.SetVehicleToWateringPlan(ctx, &sqlc.SetVehicleToWateringPlanParams{
		WateringPlanID: id,
		VehicleID:      entity.Transporter.ID,
	})
	if err != nil {
		return w.store.HandleError(err)
	}

	// link trailer to pivot table
	if entity.Trailer != nil {
		err = w.store.SetVehicleToWateringPlan(ctx, &sqlc.SetVehicleToWateringPlanParams{
			WateringPlanID: id,
			VehicleID:      entity.Trailer.ID,
		})
		if err != nil {
			return w.store.HandleError(err)
		}
	}

	return nil
}

func (w *WateringPlanRepository) setLinkedUsers(ctx context.Context, entity *entities.WateringPlan, id int32) error {
	for _, userID := range entity.UserIDs {
		err := w.store.SetUserToWateringPlan(ctx, &sqlc.SetUserToWateringPlanParams{
			WateringPlanID: id,
			UserID:         utils.UUIDToPGUUID(*userID),
		})
		if err != nil {
			return w.store.HandleError(err)
		}
	}

	return nil
}

func (w *WateringPlanRepository) setLinkedTreeClusters(ctx context.Context, entity *entities.WateringPlan, id int32) error {
	for _, tc := range entity.TreeClusters {
		err := w.store.SetTreeClusterToWateringPlan(ctx, &sqlc.SetTreeClusterToWateringPlanParams{
			WateringPlanID: id,
			TreeClusterID:  tc.ID,
		})
		if err != nil {
			return w.store.HandleError(err)
		}
	}

	return nil
}
