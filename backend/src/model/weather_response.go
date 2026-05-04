// Define estrutura/modelo para a resposta limpa de current weather.

package model

// WeatherToolResponse representa o JSON limpo que a nossa API devolve
// para a tool get_current_weather do agent.
type WeatherToolResponse struct {
	Type        string  `json:"type"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
	TempMin     float64 `json:"temp_min"`
	TempMax     float64 `json:"temp_max"`
	Humidity    int     `json:"humidity"`
	Pressure    int     `json:"pressure"`
	Description string  `json:"description"`
	WindSpeed   float64 `json:"wind_speed"`
	Clouds      int     `json:"clouds"`
	Visibility  int     `json:"visibility"`
}