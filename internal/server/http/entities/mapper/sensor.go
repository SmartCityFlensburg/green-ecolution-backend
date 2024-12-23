package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend MapSensorStatus MapLatestDataToResponse
type SensorHTTPMapper interface {
	FromResponse(src *domain.Sensor) *entities.SensorResponse
	FromWatermarkResponse(src *domain.Watermark) *entities.WatermarkResponse
}

func MapLatestDataToResponse(sensorData *domain.SensorData) *entities.SensorDataResponse {
	if sensorData.Data == nil {
		return &entities.SensorDataResponse{}
	}

	return &entities.SensorDataResponse{
		Battery:     sensorData.Data.Battery,
		Humidity:    sensorData.Data.Humidity,
		Temperature: sensorData.Data.Temperature,
		Watermarks:  mapWatermarkData(sensorData.Data.Watermarks),
	}
}

func mapWatermarkData(watermarks []domain.Watermark) []*entities.WatermarkResponse {
	responses := make([]*entities.WatermarkResponse, len(watermarks))
	for i, w := range watermarks {
		responses[i] = &entities.WatermarkResponse{
			Centibar:   w.Centibar,
			Resistance: w.Resistance,
			Depth:      w.Depth,
		}
	}
	return responses
}

func MapSensorStatus(src domain.SensorStatus) entities.SensorStatus {
	return entities.SensorStatus(src)
}
