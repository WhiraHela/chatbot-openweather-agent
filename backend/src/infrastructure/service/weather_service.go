// service de clima atual

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"web-service-gin/src/config"
	appLogger "web-service-gin/src/infrastructure/logger"
	"web-service-gin/src/infrastructure/model"
)

// WeatherError representa um erro tratado da camada de service.
// Ele carrega o status HTTP que o controller deve retornar
// e uma mensagem amigável para o usuário.
type WeatherError struct {
	StatusCode  int
	UserMessage string
	LogMessage  string
}

func (e *WeatherError) Error() string {
	return e.LogMessage
}

// GetCurrentWeather busca o clima atual de uma cidade.
// Chama OpenWeather, recebe JSON bruto, transforma em JSON mais limpo,
// e devolve para o controller retornar para a tool get_current_weather no frontend.
func GetCurrentWeather(cfg config.Config, city string) (model.WeatherToolResponse, error) {
	if cfg.OpenWeatherAPIKey == "" { // verifica se a chave da OpenWeather foi carregada
		appLogger.Critical("OPENWEATHER_API_KEY não configurada")

		return model.WeatherToolResponse{}, &WeatherError{
			StatusCode:  http.StatusInternalServerError,
			UserMessage: "OPENWEATHER_API_KEY não configurada.",
			LogMessage:  "OPENWEATHER_API_KEY não configurada",
		}
	}

	if cfg.OpenWeatherCurrentURL == "" { // verifica se a URL do endpoint de clima atual foi carregada
		appLogger.Critical("OPENWEATHER_CURRENT_URL não configurada")

		return model.WeatherToolResponse{}, &WeatherError{
			StatusCode:  http.StatusInternalServerError,
			UserMessage: "OPENWEATHER_CURRENT_URL não configurada.",
			LogMessage:  "OPENWEATHER_CURRENT_URL não configurada",
		}
	}

	appLogger.Info("Chamando OpenWeather Current para cidade: " + city)

	// codifica cidade para uso seguro dentro da URL
	encodedCity := url.QueryEscape(city)

	apiURL := fmt.Sprintf(
		"%s?q=%s&appid=%s&units=metric&lang=pt_br",
		cfg.OpenWeatherCurrentURL,
		encodedCity,
		cfg.OpenWeatherAPIKey,
	)

	client := http.Client{ // client http com timeout
		Timeout: 10 * time.Second,
	}

	// faz chamada para a OpenWeather
	resp, err := client.Get(apiURL)
	if err != nil {
		appLogger.Error("Erro ao chamar OpenWeather Current: " + err.Error())

		return model.WeatherToolResponse{}, &WeatherError{
			StatusCode:  http.StatusBadGateway,
			UserMessage: "Erro ao consultar serviço externo de clima.",
			LogMessage:  "Erro ao chamar OpenWeather Current: " + err.Error(),
		}
	}

	// fecha corpo da resposta ao final da funcao
	defer resp.Body.Close()

	// tratamento de cidade nao encontrada
	if resp.StatusCode == http.StatusNotFound {
		appLogger.Warning("Cidade não encontrada na OpenWeather: " + city)

		return model.WeatherToolResponse{}, &WeatherError{
			StatusCode:  http.StatusNotFound,
			UserMessage: "Cidade não encontrada.",
			LogMessage:  "Cidade não encontrada na OpenWeather: " + city,
		}
	}

	// tratamento de erro geral
	if resp.StatusCode != http.StatusOK {
		appLogger.Error(fmt.Sprintf("OpenWeather retornou status inesperado em /weather: %d", resp.StatusCode))

		return model.WeatherToolResponse{}, &WeatherError{
			StatusCode:  http.StatusBadGateway,
			UserMessage: "Erro retornado pela OpenWeather.",
			LogMessage:  fmt.Sprintf("OpenWeather retornou status inesperado em /weather: %d", resp.StatusCode),
		}
	}

	// struct que recebe JSON bruto da OpenWeather
	var openWeather model.OpenWeatherResponse

	// decodifica o JSON da resposta para a struct Go
	if err := json.NewDecoder(resp.Body).Decode(&openWeather); err != nil {
		appLogger.Error("Erro ao decodificar JSON da OpenWeather Current: " + err.Error())

		return model.WeatherToolResponse{}, &WeatherError{
			StatusCode:  http.StatusInternalServerError,
			UserMessage: "Erro ao processar resposta da OpenWeather.",
			LogMessage:  "Erro ao decodificar JSON da OpenWeather Current: " + err.Error(),
		}
	}

	// retorna weather como array - so acessamos o primeiro item se ele existir
	description := ""
	if len(openWeather.Weather) > 0 {
		description = openWeather.Weather[0].Description
	}

	// monta resposta limpa para o agent
	response := model.WeatherToolResponse{
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

	// responseJSON, err := json.MarshalIndent(response, "", " ")
	// 	if err != nil {
	// 		appLogger.Error("Erro ao converter resposta forecast para JSON: " + err.Error())
	// 	} else {
	// 		appLogger.Info("JSON limpo retornado para a tool: " + string(responseJSON))
	// }

	appLogger.Info("Resposta de clima atual montada com sucesso para cidade: " + response.City)

	return response, nil
}