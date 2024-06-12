package repository

import (
	"context"
	"log"
	"mqtt-golang-rainfall-prediction/models"
	"mqtt-golang-rainfall-prediction/pkg"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type RepositoryInterface interface {
	InsertData(ctx context.Context, data models.SensorData) error
	GetData(ctx context.Context) ([]models.SensorDataResponse, error)
}

type repository struct {
	influxdb influxdb2.Client
}

func NewRepository(influxdb influxdb2.Client) *repository {
	return &repository{
		influxdb: influxdb,
	}
}

func (r *repository) InsertData(ctx context.Context, data models.SensorData) error {
	writeAPI := r.influxdb.WriteAPIBlocking(os.Getenv("ORG_INFLUX"), os.Getenv("BUCKET_INFLUX"))

	point := influxdb2.NewPoint(
		"rainfall",
		map[string]string{"sensorID": "sensor1"},
		map[string]interface{}{
			"temperature":   data.Temperature,
			"humidity":      data.Humidity,
			"message":       data.Message,
			"rain_was_fall": data.RainWasFall,
			"pressure":      data.Pressure,
		},
		time.Now().UTC(),
	)

	err := writeAPI.WritePoint(ctx, point)
	if err != nil {
		return err
	}
	pkg.Broadcast <- []byte("Data berhasil diinputkan ke InfluxDB")
	return nil
}

func (r *repository) GetData(ctx context.Context) ([]models.SensorDataResponse, error) {
	queryApi := r.influxdb.QueryAPI(os.Getenv("ORG_INFLUX"))
	query := `
	from(bucket: "rainfall_data")
	|> range(start: -inf)
	|> filter(fn: (r) => r._measurement == "rainfall")
	`
	result, err := queryApi.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var resultData []models.SensorDataResponse
	for result.Next() {
		values := result.Record().Values()
		log.Println("values: ", values)

		// Create a new instance of SensorDataResponse for each record
		var data models.SensorDataResponse
		if temp, ok := values["_value"].(string); ok && values["_field"].(string) == "temperature" {
			data.Temperature = temp
		}
		if humidity, ok := values["_value"].(string); ok && values["_field"].(string) == "humidity" {
			data.Humidity = humidity
		}
		if message, ok := values["_value"].(string); ok && values["_field"].(string) == "message" {
			data.Message = message
		}
		if rainWasFall, ok := values["_value"].(string); ok && values["_field"].(string) == "rain_was_fall" {
			data.RainWasFall = rainWasFall
		}
		if pressure, ok := values["_value"].(string); ok && values["_field"].(string) == "pressure" {
			data.Pressure = pressure
		}

		// Append the populated struct to the result slice
		resultData = append(resultData, data)
	}
	return resultData, nil
}
