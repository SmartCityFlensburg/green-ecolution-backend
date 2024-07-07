package sensor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/SmartCityFlensburg/green-space-management/internal/mapper"
	"github.com/SmartCityFlensburg/green-space-management/internal/mapper/generated"
	sensorResponse "github.com/SmartCityFlensburg/green-space-management/internal/service/entities/sensor"
	"github.com/SmartCityFlensburg/green-space-management/internal/storage"
	sensorRepo "github.com/SmartCityFlensburg/green-space-management/internal/storage/entities/sensor"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttService struct {
	sensorRepo  storage.SensorRepository
	mapper      mapper.MqttMapper
	isConnected bool
}

func NewMqttService(sensorRepository storage.SensorRepository) *MqttService {
	return &MqttService{sensorRepo: sensorRepository, mapper: &generated.MqttMapperImpl{}}
}

func (s *MqttService) HandleMessage(client MQTT.Client, msg MQTT.Message) {
	jsonStr := string(msg.Payload())
	log.Printf("Received message: %s\n", jsonStr)

	var sensorData sensorResponse.MqttPayloadResponse
	if err := json.Unmarshal([]byte(jsonStr), &sensorData); err != nil {
		log.Printf("Error unmarshalling sensor data: %v\n", err)
		return
	}



	payloadEntity := s.mapper.ToEntity(
		s.mapper.FromResponse(&sensorData),
	)
	log.Printf("Mapped entity: %v\n", payloadEntity)

	entity := &sensorRepo.MqttEntity{
		Data:   *payloadEntity,
		TreeID: "6686f54fd32cf640e8ae6eb1",
	}

	if _, err := s.sensorRepo.Insert(context.Background(), entity); err != nil {
		log.Printf("Error upserting sensor data: %v\n", err)
		return
	}
}

func (s *MqttService) SetConnected(ready bool) {
	s.isConnected = ready
}

func (s *MqttService) Ready() bool {
	return s.isConnected
}
