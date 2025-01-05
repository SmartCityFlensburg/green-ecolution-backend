package mapper

import (
	"encoding/json"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	mqtt "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTimePtr
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend MapSensorStatus MapSensorData
type InternalSensorRepoMapper interface {
	// goverter:ignore LatestData
	FromSql(src *sqlc.Sensor) *entities.Sensor
	FromSqlList(src []*sqlc.Sensor) []*entities.Sensor
	// goverter:map Data | MapSensorData
	FromSqlSensorData(src *sqlc.SensorDatum) (*entities.SensorData, error)
	FromSqlSensorDataList(src []*sqlc.SensorDatum) ([]*entities.SensorData, error)
	FromDomainSensorData(src *entities.MqttPayload) *mqtt.MqttPayload
}

func MapSensorData(src []byte) (*entities.MqttPayload, error) {
	var payload entities.MqttPayload
	err := json.Unmarshal(src, &payload)
	if err != nil {
		slog.Error("Error unmarshalling sensor data", "Error", err)
		return nil, err
	}
	return &payload, nil
}

func MapSensorStatus(src sqlc.SensorStatus) entities.SensorStatus {
	return entities.SensorStatus(src)
}
