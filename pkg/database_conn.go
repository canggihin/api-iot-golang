package pkg

import (
	"context"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func ConnectInfluxDB() (influxdb2.Client, error) {
	url := os.Getenv("URL_INFLUX")
	token := os.Getenv("TOKEN_INFLUX")

	client := influxdb2.NewClient(url, token)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := client.Health(ctx)
	if err != nil {
		log.Fatal("Connection influxdb error", err)
		return nil, err
	}

	log.Println("Influxdb connected")

	return client, nil

}
