package sensor

import (
	"context"
	"errors"
	"testing"

	sensorUtils "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)


func TestSensorService_GetAll(t *testing.T) {
	t.Run("should return all sensor", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		// when
		sensorRepo.EXPECT().GetAll(context.Background()).Return(sensorUtils.TestSensorList, nil)
		sensors, err := svc.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Equal(t, sensorUtils.TestSensorList, sensors)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		sensorRepo.EXPECT().GetAll(context.Background()).Return(nil, storage.ErrSensorNotFound)
		sensors, err := svc.GetAll(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, sensors)
	})
}

func TestSensorService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete a sensor", func(t *testing.T) {
		// given
		id := int32(1)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		sensorRepo.EXPECT().GetByID(ctx, id).Return(sensorUtils.TestSensor, nil)
		treeRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		flowerbedRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		sensorRepo.EXPECT().Delete(ctx, id).Return(nil)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.NoError(t, err)
	})
	t.Run("should return error if sensor not found", func(t *testing.T) {
		// given
		id := int32(1)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		expectedErr := storage.ErrEntityNotFound
		sensorRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})

	t.Run("should return error if unlinking sensor ID on tree fails", func(t *testing.T) {
		// given
		id := int32(1)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		expectedErr := errors.New("failed to unlink")

		sensorRepo.EXPECT().GetByID(ctx, id).Return(sensorUtils.TestSensor, nil)
		treeRepo.EXPECT().UnlinkSensorID(ctx, id).Return(expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})

	t.Run("should return error if unlinking sensor ID on flowerbed fails", func(t *testing.T) {
		// given
		id := int32(1)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)
		expectedErr := errors.New("failed to unlink")

		sensorRepo.EXPECT().GetByID(ctx, id).Return(sensorUtils.TestSensor, nil)
		treeRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		flowerbedRepo.EXPECT().UnlinkSensorID(ctx, id).Return(expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})

	t.Run("should return error if deleting sensor fails", func(t *testing.T) {
		// given
		id := int32(1)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)
		expectedErr := errors.New("failed to delete")

		sensorRepo.EXPECT().GetByID(ctx, id).Return(sensorUtils.TestSensor, nil)
		treeRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		flowerbedRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		sensorRepo.EXPECT().Delete(ctx, id).Return(expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})
}

func TestReady(t *testing.T) {
	t.Run("should return true if the service is ready", func(t *testing.T) {
		// given
		repo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(repo, treeRepo, flowerbedRepo)

		// when
		ready := svc.Ready()

		// then
		assert.True(t, ready)
	})

	t.Run("should return false if the service is not ready", func(t *testing.T) {
		// give
		svc := NewSensorService(nil, nil, nil)

		// when
		ready := svc.Ready()

		// then
		assert.False(t, ready)
	})
}
