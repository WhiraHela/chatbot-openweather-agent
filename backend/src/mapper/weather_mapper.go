package mapper

import "web-service-gin/src/model"

// ToWeatherToolResponse transforma o JSON bruto da OpenWeather
// em uma resposta limpa para a tool get_current_weather do agent.
func ToWeatherToolResponse(openWeather model.OpenWeatherResponse) model.WeatherToolResponse {
	// retorna weather como array - so acessamos o primeiro item se ele existir
	description := ""
	if len(openWeather.Weather) > 0 {
		description = openWeather.Weather[0].Description
	}

	// monta resposta limpa para o agent
	return model.WeatherToolResponse{
		Type:        "current_weather",
		City:        openWeather.Name,
		Country:     openWeather.Sys.Country,
		Temperature: openWeather.Main.Temp,
		FeelsLike:   openWeather.Main.FeelsLike,
		TempMin:     openWeather.Main.TempMin,
		TempMax:     openWeather.Main.TempMax,
		Humidity:    openWeather.Main.Humidity,
		Pressure:    openWeather.Main.Pressure,
		Description: description,
		WindSpeed:   openWeather.Wind.Speed,
		Clouds:      openWeather.Clouds.All,
		Visibility:  openWeather.Visibility,
	}
}

// ToForecastToolResponse transforma o JSON bruto da OpenWeather
// em uma resposta limpa para a tool get_weather_forecast do agent.
func ToForecastToolResponse(forecast model.OpenWeatherForecastResponse, days int) model.ForecastToolResponse {
	// calcula quantidade de registros por dia (8)
	limit := days * 8
	if limit > len(forecast.List) {
		limit = len(forecast.List)
	}

	// cria slice vazio de ForecastItem com capacidade inicial igual ao limite
	// para evitar realocacoes desnecessarias
	items := make([]model.ForecastItem, 0, limit)

	// percorre apenas pelos registros necessarios, de acordo com numero de dias
	for i := 0; i < limit; i++ {
		item := forecast.List[i]

		description := ""
		if len(item.Weather) > 0 {
			description = item.Weather[0].Description
		}

		// transforma o item bruto da OpenWeather em uma estrutura mais simples
		// para o agent e frontend consumirem
		forecastItem := model.ForecastItem{
			DateTime:         item.DateText,
			Temperature:      item.Main.Temp,
			FeelsLike:        item.Main.FeelsLike,
			TempMin:          item.Main.TempMin,
			TempMax:          item.Main.TempMax,
			Humidity:         item.Main.Humidity,
			Description:      description,
			WindSpeed:        item.Wind.Speed,
			RainProbability:  item.Pop,
			RainVolume3h:     item.Rain.ThreeHours,
			Clouds:           item.Clouds.All,
		}

		items = append(items, forecastItem)
	}

	// monta resposta final limpa -> resposta que a tool do agent vai receber
	return model.ForecastToolResponse{
		Type:          "weather_forecast",
		City:          forecast.City.Name,
		Country:       forecast.City.Country,
		DaysRequested: days,
		Items:         items,
	}
}