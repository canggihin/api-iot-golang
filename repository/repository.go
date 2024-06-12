package repository

import (
	"context"
	"encoding/json"
	"log"
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
			"data": data},
		time.Now().UTC(),
	)

	err := writeAPI.WritePoint(ctx, point)
	if err != nil {
		return err
	}
	pkg.Broadcast <- []byte("data berhasil diinputkan ke influxdb")
	return nil
}

func (r *repository) GetData(ctx context.Context) ([]models.SensorData, error) {
	queryApi := r.influxdb.QueryAPI(os.Getenv("ORG_INFLUX"))
	query := `
	from(bucket: "rainfall_data")
	|> range(start: -inf)
	|> filter(fn: (r) => r._measurement == "rainfall")
	|> filter(fn: (r) => r._field == "data")
	|> keep(columns: ["_value"])
	`
	result, err := queryApi.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var resultData []models.SensorData
	for result.Next() {
		var data models.SensorData
		record := result.Record().Value().(string)
		log.Println("record: ", record)
		if err := json.Unmarshal([]byte(record), &data); err != nil {
			log.Println("Error unmarshal data: ", err)
			continue
		}
		resultData = append(resultData, data)
	}
	return resultData, nil
}
