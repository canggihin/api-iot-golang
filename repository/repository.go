package repository

import (
	"context"
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
	|> keep(columns: ["_time", "temperature", "humidity", "message", "rain_was_fall", "pressure"])
	`
	result, err := queryApi.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var resultData []models.SensorDataResponse
	for result.Next() {
		var data models.SensorDataResponse
		values := result.Record().Values()
		data.Temperature = values["temperature"].(string)
		data.Humidity = values["humidity"].(string)
		data.Message = values["message"].(string)
		data.RainWasFall = values["rain_was_fall"].(string)
		data.Pressure = values["pressure"].(string)

		resultData = append(resultData, data)
	}
	return resultData, nil
}
