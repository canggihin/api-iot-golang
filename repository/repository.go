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
	GetData(ctx context.Context) ([]models.SensorData, error)
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

func (r *repository) GetData(ctx context.Context) ([]models.SensorData, error) {
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

	dataMap := make(map[time.Time]*models.SensorData)
	var resultData []models.SensorData
	for result.Next() {
		values := result.Record().Values()
		timestamp := result.Record().Time() // Get the timestamp of the record
		if dataMap[timestamp] == nil {
			dataMap[timestamp] = &models.SensorData{} // Initialize if not already
		}
		data := dataMap[timestamp]

		// Assign data based on field type
		switch values["_field"].(string) {
		case "temperature":
			if temp, ok := values["_value"].(float64); ok {
				data.Temperature = temp
			}
		case "humidity":
			if humidity, ok := values["_value"].(float64); ok {
				data.Humidity = humidity
			}
		case "message":
			if message, ok := values["_value"].(string); ok {
				data.Message = message
			}
		case "rain_was_fall":
			if rainWasFall, ok := values["_value"].(float64); ok {
				data.RainWasFall = rainWasFall
			}
		case "pressure":
			if pressure, ok := values["_value"].(float64); ok {
				data.Pressure = pressure
			}
		}
	}

	// Transfer from map to slice
	for _, d := range dataMap {
		resultData = append(resultData, *d)
	}

	return resultData, nil
}
