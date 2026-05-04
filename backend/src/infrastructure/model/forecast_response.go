// Define estrutura/modelo para a resposta limpa de forecast weather.

package model

// ForecastToolResponse representa o JSON limpo que a nossa API devolve
// para a tool get_weather_forecast do agent.
type ForecastToolResponse struct {
	Type          string         `json:"type"`           // tipo da resposta
	City          string         `json:"city"`           // nome da cidade retornada pelo OpenWeather
	Country       string         `json:"country"`        // país
	DaysRequested int            `json:"days_requested"` // quantidade de dias solicitados
	Items         []ForecastItem `json:"items"`          // lista de previsões em intervalos de 3 horas
}

// ForecastItem representa cada item individual da previsão.
type ForecastItem struct {
	DateTime         string  `json:"date_time"`
	Temperature     float64 `json:"temperature"`
	FeelsLike       float64 `json:"feels_like"`
	TempMin         float64 `json:"temp_min"`
	TempMax         float64 `json:"temp_max"`
	Humidity        int     `json:"humidity"`
	Description     string  `json:"description"`
	WindSpeed       float64 `json:"wind_speed"`
	RainProbability float64 `json:"rain_probability"`
	RainVolume3h    float64 `json:"rain_volume_3h"`
	Clouds          int     `json:"clouds"`
}