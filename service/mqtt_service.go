package service

import (
	"context"
	"encoding/json"
	"log"
	"mqtt-golang-rainfall-prediction/models"
	"mqtt-golang-rainfall-prediction/pkg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (s *service) MessageMqttHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("TOPIC: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
	var data models.SystemInfo
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Println("Error unmarshal data: ", err)
		return
	}
	log.Println("Data: ", data)
	_, err := s.ProsesMessage(context.Background(), data)
	if err != nil {
		log.Println("Error proses message: ", err)
		return
	}
}

func (s *service) ProsesMessage(ctx context.Context, data models.SystemInfo) (models.SystemInfo, error) {
	if data == (models.SystemInfo{}) {
		return models.SystemInfo{}, nil
	}

	pkg.BroadcastToSystems([]byte(`{"status": 1, "message": "Data accepted"}`))
	return data, nil
}
