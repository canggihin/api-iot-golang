package mqtt_pkg

import (
	"fmt"
	"log"
	"mqtt-golang-rainfall-prediction/repository"
	"mqtt-golang-rainfall-prediction/service"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func ConnectMqtt(influxdb influxdb2.Client) mqtt.Client {
	broker := os.Getenv("MQTT_BROKER")
	port := os.Getenv("MQTT_PORT")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", broker, port))
	opts.SetClientID("mqtt-golang-rainfall-prediction")

	repo := repository.NewRepository(influxdb)
	service := service.NewService(repo)
	opts.SetDefaultPublishHandler(service.MessageMqttHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetCleanSession(false)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func SubscribeMqtt(mqtt mqtt.Client, topic string, influxdb influxdb2.Client) bool {
	repo := repository.NewRepository(influxdb)
	service := service.NewService(repo)
	topics := fmt.Sprintf("rainfall/%s", topic)
	token := mqtt.Subscribe(topics, 1, service.MessageMqttHandler)
	success := token.Wait()
	if !success {
		return false
	}
	log.Printf("Subscribed to topic: %s\n", topics)
	return true
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}
