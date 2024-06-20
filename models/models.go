package models

type SensorData struct {
	Temperature   float64 `json:"temperature"`
	Humidity      float64 `json:"humidity"`
	Message       string  `json:"message"`
	RainWasFall   float64 `json:"rain_was_fall"`
	Pressure      float64 `json:"pressure"`
	FormattedTime string  `json:"formattedTime"`
}

type SensorDataByDay struct {
	Temperature   float64 `json:"temperature"`
	Humidity      float64 `json:"humidity"`
	RainWasFall   float64 `json:"rain_was_fall"`
	Pressure      float64 `json:"pressure"`
	FormattedTime string  `json:"formattedTime"`
}

type SystemInfo struct {
	TotalSensor  int    `json:"total_sensor"`
	RamConsume   string `json:"ram_consume"`
	CpuConsume   string `json:"cpu_consume"`
	DHTSensor    int    `json:"dht_sensor"`
	BMP180Sensor int    `json:"bmp180_sensor"`
	RainSensor   int    `json:"rain_sensor"`
	BatteryLevel int    `json:"battery_level"`
}
