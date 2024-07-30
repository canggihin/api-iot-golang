package main

import (
	"fmt"
	"mqtt-golang-rainfall-prediction/handlers"
	"mqtt-golang-rainfall-prediction/middleware"
	"mqtt-golang-rainfall-prediction/pkg"
	mqtt_pkg "mqtt-golang-rainfall-prediction/pkg/mqtt"
	"mqtt-golang-rainfall-prediction/repository"
	"mqtt-golang-rainfall-prediction/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

func router(r *gin.Engine, influxdb influxdb2.Client) {

	repo := repository.NewRepository(influxdb)
	service := service.NewService(repo)
	handlers := handlers.NewHandler(service)

	r.POST("/data", middleware.AuthMiddleware("user", "admin", "superadmin"), handlers.CreateData)
	r.GET("/ws/sensor", handlers.HandleWsSensor)
	r.GET("/ws/system", handlers.HandleWsSystem)
	r.GET("/data", handlers.GetData)
	r.GET("/reportday", middleware.AuthMiddleware("user", "admin"), handlers.GetDataByDay)
	r.POST("/connect", handlers.SuccessConnectedDevice)
	r.POST("/sysinfo", handlers.GetSystemInfo)
}

func main() {
	_ = godotenv.Load(".env")
	influxdb, err := pkg.ConnectInfluxDB()

	if err != nil {
		fmt.Println("Error connecting to InfluxDB")
		return
	}

	clientMqtt := mqtt_pkg.ConnectMqtt(influxdb)
	mqtt_pkg.SubscribeMqtt(clientMqtt, "sensor", influxdb)
	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true

	r := gin.Default()

	r.Use(cors.New(configCors))
	router(r, influxdb)

	r.Run(":8089")
	fmt.Println("HTTP server is running on http://localhost:8089")
}
