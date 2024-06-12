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

func handleMessages() {
	for {
		msg := <-pkg.Broadcast
		for client := range pkg.Clients {
			err := client.WriteMessage(1, msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(pkg.Clients, client)
			}
		}
	}
}

func router(r *gin.Engine, influxdb influxdb2.Client) {

	repo := repository.NewRepository(influxdb)
	service := service.NewService(repo)
	handlers := handlers.NewHandler(service)

	r.POST("/data", handlers.CreateData)
	r.GET("/ws", handlers.HandleConnectionWs)
	r.GET("/data", handlers.GetData)
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

	go handleMessages()
	r.Run(":8089")
	fmt.Println("HTTP server is running on http://localhost:8089")
}
