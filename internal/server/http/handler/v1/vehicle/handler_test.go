package vehicle

import (
	"context"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestGetAllVehicles(t *testing.T) {
	t.Run("should return all vehicles successfully", func(t *testing.T) {
		app := fiber.New()
		mockVehicleService := serviceMock.NewMockVehicleService(t)
		handler := GetAllVehicles(mockVehicleService)
		app.Get("/v1/vehicle", handler)

		mockVehicleService.EXPECT().GetAll(
			mock.Anything,
		).Return(TestVehicles, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/vehicle", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.VehicleListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Data))
		assert.Equal(t, TestVehicles[0].ID, response.Data[0].ID)

		mockVehicleService.AssertExpectations(t)
	})

	t.Run("should return an empty list when no vehicles are available", func(t *testing.T) {
		app := fiber.New()
		mockVehicleService := serviceMock.NewMockVehicleService(t)
		handler := GetAllVehicles(mockVehicleService)
		app.Get("/v1/vehicle", handler)

		mockVehicleService.EXPECT().GetAll(
			mock.Anything,
		).Return([]*entities.Vehicle{}, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/vehicle", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.VehicleListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response.Data))

		mockVehicleService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error when service fails", func(t *testing.T) {
		app := fiber.New()
		mockVehicleService := serviceMock.NewMockVehicleService(t)
		handler := GetAllVehicles(mockVehicleService)
		app.Get("/v1/vehicle", handler)

		mockVehicleService.EXPECT().GetAll(
			mock.Anything,
			).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/vehicle", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockVehicleService.AssertExpectations(t)
	})
}

func TestGetVehicleByID(t *testing.T) {
	t.Run("should return vehicle successfully", func(t *testing.T) {
		app := fiber.New()
		mockVehicleService := serviceMock.NewMockVehicleService(t)
		handler := GetVehicleByID(mockVehicleService)
		app.Get("/v1/vehicle/:id", handler)

		mockVehicleService.EXPECT().GetByID(
			mock.Anything,
			int32(1),
		).Return(TestVehicle, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/vehicle/1", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.TreeClusterResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, TestVehicle.ID, response.ID)

		mockVehicleService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid vehicle ID", func(t *testing.T) {
		app := fiber.New()
		mockVehicleService := serviceMock.NewMockVehicleService(t)
		handler := GetVehicleByID(mockVehicleService)
		app.Get("/v1/vehicle/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/vehicle/invalid-id", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 Not Found if vehicle does not exist", func(t *testing.T) {
		app := fiber.New()
		mockVehicleService := serviceMock.NewMockVehicleService(t)
		handler := GetVehicleByID(mockVehicleService)
		app.Get("/v1/vehicle/:id", handler)

		mockVehicleService.EXPECT().GetByID(
			mock.Anything,
			int32(999),
		).Return(nil, storage.ErrVehicleNotFound)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/vehicle/999", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockVehicleService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error for service failure", func(t *testing.T) {
		app := fiber.New()
		mockVehicleService := serviceMock.NewMockVehicleService(t)
		handler := GetVehicleByID(mockVehicleService)
		app.Get("/v1/vehicle/:id", handler)

		mockVehicleService.EXPECT().GetByID(
			mock.Anything,
			int32(1),
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/vehicle/1", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockVehicleService.AssertExpectations(t)
	})
}