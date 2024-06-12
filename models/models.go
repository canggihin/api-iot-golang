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
