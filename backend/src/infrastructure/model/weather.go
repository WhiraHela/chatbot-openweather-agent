// Define estrutura/modelo para current weather

package model

//JSON limpo que a nossa API devolve para a tool do agent
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

//JSON bruto da Openweather
type OpenWeatherResponse struct {
	Name string `json:"name"`

	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`

	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
	} `json:"main"`

	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`

	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`

	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`

	Visibility int `json:"visibility"`
}