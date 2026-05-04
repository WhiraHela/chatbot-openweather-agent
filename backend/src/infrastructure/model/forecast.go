// Define estrutura/modelo para forecast weather

package model


//JSON limpo que a nossa API devolve para a tool do agent
type ForecastToolResponse struct {
	Type          string         `json:"type"` 				//tipo daresposta
	City          string         `json:"city"` 				// nome da cidade retornada pelo OpenWeather
	Country       string         `json:"country"` 			//país
	DaysRequested int            `json:"days_requested"` 	//quantidades de dias solicitados
	Items         []ForecastItem `json:"items"` 			// lista de previsoes em intervalos de 3 horas
}

// Representa cada item individual da previsao
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


//JSON bruto da Openweather
type OpenWeatherForecastResponse struct {
	Cod     string `json:"cod"`
	Message int    `json:"message"`
	Cnt     int    `json:"cnt"`

	List []struct {
		Dt int64 `json:"dt"`

		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`

		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`

		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`

		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
			Gust  float64 `json:"gust"`
		} `json:"wind"`

		Pop float64 `json:"pop"`

		Rain struct {
			ThreeHours float64 `json:"3h"`
		} `json:"rain"`

		DateText string `json:"dt_txt"`
	} `json:"list"`

	City struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"city"`
}