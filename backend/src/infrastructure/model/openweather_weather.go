// Define estrutura/modelo para o JSON bruto de current weather da OpenWeather.

package model

// OpenWeatherResponse representa o JSON bruto retornado pela OpenWeather
// no endpoint de clima atual.
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