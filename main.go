package main

import (
	"fmt"
	"mqtt-golang-rainfall-prediction/handlers"
	"mqtt-golang-rainfall-prediction/pkg"
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

	r.POST("/data", handlers.CreateData)
	r.GET("/ws/sensor", handlers.HandleWsSensor)
	r.GET("/ws/system", handlers.HandleWsSystem)
	r.GET("/data", handlers.GetData)
	r.GET("/reportday", handlers.GetDataByDay)
	r.POST("/connect", handlers.SuccessConnectedDevice)
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

	r := gin.Default()

	r.Use(cors.New(configCors))
	router(r, influxdb)

	r.Run(":8089")
	fmt.Println("HTTP server is running on http://localhost:8089")
}
