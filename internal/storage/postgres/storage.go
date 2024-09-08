package postgres

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/flowerbed"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image"
	mapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/vehicle"
	"github.com/jackc/pgx/v5"
)

func NewRepository(conn *pgx.Conn) *storage.Repository {
	s := store.NewStore(conn)

	treeMappers := tree.NewTreeRepositoryMappers(
		&mapper.InternalTreeRepoMapperImpl{},
		&mapper.InternalImageRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
		&mapper.InternalTreeClusterRepoMapperImpl{},
	)
	treeRepo := tree.NewTreeRepository(s, treeMappers)

	tcMappers := treecluster.NewTreeClusterRepositoryMappers(
		&mapper.InternalTreeClusterRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
	)
	treeClusterRepo := treecluster.NewTreeClusterRepository(s, tcMappers)

	imageMappers := image.NewImageRepositoryMappers(
		&mapper.InternalImageRepoMapperImpl{},
	)
	imageRepo := image.NewImageRepository(s, imageMappers)

	vehicleMappers := vehicle.NewVehicleRepositoryMappers(
		&mapper.InternalVehicleRepoMapperImpl{},
	)
	vehicleRepo := vehicle.NewVehicleRepository(s, vehicleMappers)

	sensorMappers := sensor.NewSensorRepositoryMappers(
		&mapper.InternalSensorRepoMapperImpl{},
	)
	sensorRepo := sensor.NewSensorRepository(s, sensorMappers)

	flowMappers := flowerbed.NewFlowerbedMappers(
		&mapper.InternalFlowerbedRepoMapperImpl{},
		&mapper.InternalImageRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
	)
	flowerbedRepo := flowerbed.NewFlowerbedRepository(s, flowMappers)

	return &storage.Repository{
		Tree:        treeRepo,
		TreeCluster: treeClusterRepo,
		Image:       imageRepo,
		Vehicle:     vehicleRepo,
		Sensor:      sensorRepo,
		Flowerbed:   flowerbedRepo,
	}
}
