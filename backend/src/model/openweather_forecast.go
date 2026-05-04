// Define estrutura/modelo para o JSON bruto de forecast weather da OpenWeather.

package model

// OpenWeatherForecastResponse representa o JSON bruto retornado pela OpenWeather
// no endpoint de previsão do tempo.
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