package openrouteservice

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"

	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/openrouteservice/ors"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/vroom"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

// validate is RouteRepo implements storage.RoutingRepository
var _ storage.RoutingRepository = (*RouteRepo)(nil)

type RouteRepoConfig struct {
	routing config.RoutingConfig
}

type RouteRepo struct {
	vroom vroom.VroomClient
	ors   ors.OrsClient
	cfg   *RouteRepoConfig
}

func NewRouteRepo(cfg *RouteRepoConfig) (*RouteRepo, error) {
	vroomURL, err := url.Parse(cfg.routing.Ors.Optimization.Vroom.Host)
	if err != nil {
		return nil, err
	}
	orsURL, err := url.Parse(cfg.routing.Ors.Host)
	if err != nil {
		return nil, err
	}

	vroomClient := vroom.NewVroomClient(
		vroom.WithHostURL(vroomURL),
		vroom.WithStartPoint(cfg.routing.StartPoint),
		vroom.WithEndPoint(cfg.routing.EndPoint),
		vroom.WithWateringPoint(cfg.routing.WateringPoint),
	)
	orsClient := ors.NewOrsClient(
		ors.WithHostURL(orsURL),
	)

	return &RouteRepo{
		vroom: vroomClient,
		ors:   orsClient,
		cfg:   cfg,
	}, nil
}

func (r *RouteRepo) GenerateRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (*entities.GeoJSON, error) {
	orsProfile, err := r.toOrsVehicleType(vehicle.Type)
	if err != nil {
		return nil, err
	}

	_, orsRoute, err := r.prepareOrsRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	entity, err := r.ors.DirectionsGeoJSON(ctx, orsProfile, orsRoute)
	if err != nil {
		return nil, err
	}

	metadata, err := routing.ConvertLocations(&r.cfg.routing)
	if err != nil {
		return nil, err
	}

	entity.Metadata = *metadata

	return entity, nil
}

func (r *RouteRepo) GenerateRawGpxRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (io.ReadCloser, error) {
	orsProfile, err := r.toOrsVehicleType(vehicle.Type)
	if err != nil {
		return nil, err
	}

	_, orsRoute, err := r.prepareOrsRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	return r.ors.DirectionsRawGpx(ctx, orsProfile, orsRoute)
}

func (r *RouteRepo) GenerateRouteInformation(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (*entities.RouteMetadata, error) {
	orsProfile, err := r.toOrsVehicleType(vehicle.Type)
	if err != nil {
		return nil, err
	}

	optimizedRoutes, route, err := r.prepareOrsRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	var refillCount int
	if len(optimizedRoutes.Routes) > 0 {
		oRoute := optimizedRoutes.Routes[0]
		reducedSteps := utils.Reduce(oRoute.Steps, r.reduceSteps, make([]*vroom.VroomRouteStep, 0, len(oRoute.Steps)))
		refillCount = len(utils.Filter(reducedSteps, func(step *vroom.VroomRouteStep) bool {
			return step.Type == "pickup"
		}))
	}

	rawDirections, err := r.ors.DirectionsJSON(ctx, orsProfile, route)
	if err != nil {
		return nil, err
	}

	var distance, duration float64
	if len(rawDirections.Routes) > 0 {
		distance = rawDirections.Routes[0].Summary.Distance
		duration = rawDirections.Routes[0].Summary.Duration
	}

	return &entities.RouteMetadata{
		Refills:  refillCount,
		Distance: distance,
		Time:     duration,
	}, nil

}

func (r *RouteRepo) prepareOrsRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (optimized *vroom.VroomResponse, routes *ors.OrsDirectionRequest, err error) {
	optimizedRoutes, err := r.vroom.OptimizeRoute(ctx, vehicle, clusters)
	if err != nil {
		slog.Error("failed to optimize route", "error", err)
		return nil, nil, err
	}

	// currently handle only the first route
	if len(optimizedRoutes.Routes) == 0 {
		slog.Error("there are no routes in vroom response", "routes", optimizedRoutes)
		return nil, nil, errors.New("empty routes")
	}
	oRoute := optimizedRoutes.Routes[0]
	reducedSteps := utils.Reduce(oRoute.Steps, r.reduceSteps, make([]*vroom.VroomRouteStep, 0, len(oRoute.Steps)))
	orsLocation := utils.Reduce(reducedSteps, func(acc [][]float64, current *vroom.VroomRouteStep) [][]float64 {
		return append(acc, current.Location)
	}, make([][]float64, 0, len(reducedSteps)))

	return optimizedRoutes, &ors.OrsDirectionRequest{
		Coordinates: orsLocation,
		Units:       "m",
		Language:    "de-de",
	}, nil
}

func (r *RouteRepo) toOrsVehicleType(vehicle entities.VehicleType) (string, error) {
	if vehicle == entities.VehicleTypeUnknown {
		return "", storage.ErrUnknownVehicleType
	}

	return "driving-car", nil // ORS doesn't support dynamic routing over api call
}

// Reduce multiple pickups to one
// "start" -> "pickup" -> "pickup" -> "delivery" => "start" -> "pickup" -> "delivery"
func (r *RouteRepo) reduceSteps(acc []*vroom.VroomRouteStep, current vroom.VroomRouteStep) []*vroom.VroomRouteStep {
	if len(acc) == 0 {
		return append(acc, &current)
	}

	prev := acc[len(acc)-1]
	if prev.Type != "pickup" {
		return append(acc, &current)
	}

	if current.Type != "pickup" {
		return append(acc, &current)
	}

	prev.Load = current.Load
	return acc
}
