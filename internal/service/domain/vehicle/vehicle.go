package vehicle

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/vehicle"
)

type VehicleService struct {
	vehicleRepo storage.VehicleRepository
	validator   *validator.Validate
}

func NewVehicleService(vehicleRepository storage.VehicleRepository) service.VehicleService {
	return &VehicleService{
		vehicleRepo: vehicleRepository,
		validator:   validator.New(),
	}
}

func (v *VehicleService) GetAll(ctx context.Context) ([]*entities.Vehicle, error) {
	vehicles, err := v.vehicleRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return vehicles, nil
}

func (v *VehicleService) GetByID(ctx context.Context, id int32) (*entities.Vehicle, error) {
	got, err := v.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return got, nil
}

func (v *VehicleService) GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error) {
	got, err := v.vehicleRepo.GetByPlate(ctx, plate)
	if err != nil {
		return nil, handleError(err)
	}

	return got, nil
}

func (v *VehicleService) Create(ctx context.Context, vh *entities.VehicleCreate) (*entities.Vehicle, error) {
	if err := v.validator.Struct(vh); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	if isTaken, err := v.isVehicleNumberPlateTaken(ctx, vh.NumberPlate); err != nil {
		return nil, err
	} else if isTaken {
		return nil, service.NewError(service.BadRequest, errors.New("Number plate is already taken").Error())
	}

	created, err := v.vehicleRepo.Create(ctx,
		vehicle.WithNumberPlate(vh.NumberPlate),
		vehicle.WithDescription(vh.Description),
		vehicle.WithWaterCapacity(vh.WaterCapacity),
		vehicle.WithVehicleStatus(vh.Status),
		vehicle.WithVehicleType(vh.Type),
		vehicle.WithHeight(vh.Height),
		vehicle.WithLength(vh.Length),
		vehicle.WithModel(vh.Model),
		vehicle.WithWidth(vh.Width),
		vehicle.WithDrivingLicense(vh.DrivingLicense),
	)

	if err != nil {
		return nil, handleError(err)
	}

	return created, nil
}

func (v *VehicleService) Update(ctx context.Context, id int32, vh *entities.VehicleUpdate) (*entities.Vehicle, error) {
	if err := v.validator.Struct(vh); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	_, err := v.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	if isTaken, err := v.isVehicleNumberPlateTaken(ctx, vh.NumberPlate); err != nil {
		return nil, err
	} else if isTaken {
		return nil, service.NewError(service.BadRequest, errors.New("Number plate is already taken").Error())
	}

	updated, err := v.vehicleRepo.Update(ctx,
		id,
		vehicle.WithNumberPlate(vh.NumberPlate),
		vehicle.WithDescription(vh.Description),
		vehicle.WithWaterCapacity(vh.WaterCapacity),
		vehicle.WithVehicleStatus(vh.Status),
		vehicle.WithVehicleType(vh.Type),
		vehicle.WithHeight(vh.Height),
		vehicle.WithLength(vh.Length),
		vehicle.WithModel(vh.Model),
		vehicle.WithWidth(vh.Width),
		vehicle.WithDrivingLicense(vh.DrivingLicense),
	)
	if err != nil {
		return nil, handleError(err)
	}

	return updated, nil
}

func (v *VehicleService) Delete(ctx context.Context, id int32) error {
	_, err := v.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	err = v.vehicleRepo.Delete(ctx, id)
	if err != nil {
		return handleError(err)
	}

	return nil
}

func (v *VehicleService) Ready() bool {
	return v.vehicleRepo != nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, storage.ErrVehicleNotFound.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (v *VehicleService) isVehicleNumberPlateTaken(ctx context.Context, plate string) (bool, error) {
	existingVehicle, err := v.vehicleRepo.GetByPlate(ctx, plate)
	if err != nil && !errors.Is(err, storage.ErrEntityNotFound) {
		return false, handleError(err)
	}
	return existingVehicle != nil, nil
}
