package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTimePtr
// goverter:extend MapDrivingLicense MapVehicleStatus MapVehicleType
type InternalVehicleRepoMapper interface {
	FromSql(src *sqlc.Vehicle) *entities.Vehicle
	FromSqlList(src []*sqlc.Vehicle) []*entities.Vehicle
}

func MapVehicleStatus(vehicleStatus sqlc.VehicleStatus) entities.VehicleStatus {
	return entities.VehicleStatus(vehicleStatus)
}

func MapVehicleType(vehicleType sqlc.VehicleType) entities.VehicleType {
	return entities.VehicleType(vehicleType)
}

func MapDrivingLicense(drivingLicense sqlc.DrivingLicense) entities.DrivingLicense {
	return entities.DrivingLicense(drivingLicense)
}
