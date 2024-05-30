package main

import (
	"fmt"
	"mqtt-golang-rainfall-prediction/handlers"
	"mqtt-golang-rainfall-prediction/pkg"
	"mqtt-golang-rainfall-prediction/repository"
	"mqtt-golang-rainfall-prediction/service"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

func router(r *gin.Engine, influxdb influxdb2.Client) {

	repo := repository.NewRepository(influxdb)
	service := service.NewService(repo)
	handlers := handlers.NewHandler(service)

	r.POST("/data", handlers.CreateData)
}
func ConnectMqttClient() mqtt.Client {
	broker := "36.92.168.180"
	port := 7483

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("temp-humd-mqtt-client_1")

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client

}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func publish(client mqtt.Client, topics string) {
	text := map[string]interface{}{
		"sn":      "LWD-A10TH05240001",
		"counter": "4043",
		"sensor": map[string]string{
			"air_temperature": "25.53",
			"air_humidity":    "69.98",
		},
		"config": map[string]string{
			"ssid":     "DBT",
			"password": "telkom2021",
			"interval": "5.00",
		},
		"rssi": "-65",
	}
	token := client.Publish(topics, 1, false, text)
	token.Wait()
	fmt.Println("Published")
}

func main() {
	_ = godotenv.Load(".env")
	influxdb, err := pkg.ConnectInfluxDB()

	if err != nil {
		fmt.Println("Error connecting to InfluxDB")
		return
	}

	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	client := ConnectMqttClient()
	for i := 0; i < 10; i++ {
		publish(client, "manufacture/sensor/temphum/LN-DW-04")
		time.Sleep(5 * time.Second)
	}

	r := gin.Default()

	r.Use(cors.New(configCors))
	router(r, influxdb)

	r.Run(":8089")
	fmt.Println("HTTP server is running on http://localhost:8089")
}
