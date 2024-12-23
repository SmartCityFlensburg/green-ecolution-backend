package wateringplan

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type WateringPlanRepository struct {
	store *store.Store
	WateringPlanMappers
}

type WateringPlanMappers struct {
	mapper        mapper.InternalWateringPlanRepoMapper
	vehicleMapper mapper.InternalVehicleRepoMapper
	clusterMapper mapper.InternalTreeClusterRepoMapper
}

func NewWateringPlanRepositoryMappers(
	wMapper mapper.InternalWateringPlanRepoMapper,
	vMapper mapper.InternalVehicleRepoMapper,
	tcMapper mapper.InternalTreeClusterRepoMapper,
) WateringPlanMappers {
	return WateringPlanMappers{
		mapper:        wMapper,
		vehicleMapper: vMapper,
		clusterMapper: tcMapper,
	}
}

func NewWateringPlanRepository(s *store.Store, mappers WateringPlanMappers) storage.WateringPlanRepository {
	return &WateringPlanRepository{
		store:               s,
		WateringPlanMappers: mappers,
	}
}

func (w *WateringPlanRepository) Delete(ctx context.Context, id int32) error {
	_, err := w.store.DeleteWateringPlan(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
