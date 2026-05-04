package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"web-service-gin/src/config"
	appLogger "web-service-gin/src/logger"
	"web-service-gin/src/mapper"
	"web-service-gin/src/model"
)

// ForecastError representa um erro tratado da camada de service.
// Ele carrega o status HTTP que o controller deve retornar e uma mensagem amigável para o usuário.
type ForecastError struct {
	StatusCode  int
	UserMessage string
	LogMessage  string
}

func (e *ForecastError) Error() string {
	return e.LogMessage
}

// GetForecast concentra a regra de negócio da previsão do tempo.
// Ele chama a OpenWeather, decodifica o JSON bruto e monta uma resposta limpa para o agent/frontend.
func GetForecast(cfg config.Config, city string, days int) (model.ForecastToolResponse, error) {
	// verificar se o OpenWeather foi carregada corretamente
	if cfg.OpenWeatherAPIKey == "" {
		appLogger.Critical("OPENWEATHER_API_KEY não configurada")

		return model.ForecastToolResponse{}, &ForecastError{
			StatusCode:  http.StatusInternalServerError, // codigo 500
			UserMessage: "OPENWEATHER_API_KEY não configurada.",
			LogMessage:  "OPENWEATHER_API_KEY não configurada",
		}
	}

	// verifica se a URL do endpoint forecast foi configurada
	if cfg.OpenWeatherForecastURL == "" {
		appLogger.Critical("OPENWEATHER_FORECAST_URL não configurada")

		return model.ForecastToolResponse{}, &ForecastError{
			StatusCode:  http.StatusInternalServerError, // codigo 500
			UserMessage: "OPENWEATHER_FORECAST_URL não configurada.",
			LogMessage:  "OPENWEATHER_FORECAST_URL não configurada",
		}
	}

	// trata caracteres especiais para montar a URL da API externa corretamente
	encodedCity := url.QueryEscape(city)

	// service -> puxar de formatter (classe weatherFormatter -> func forecast_formatter)
	// monta a URL final do OpenWeather
	apiURL := fmt.Sprintf(
		"%s?q=%s&appid=%s&units=metric&lang=pt_br",
		cfg.OpenWeatherForecastURL,
		encodedCity,
		cfg.OpenWeatherAPIKey,
	)

	// cria cliente http com timeout
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	appLogger.Info(fmt.Sprintf("Chamando OpenWeather Forecast para cidade: %s, dias: %d", city, days))

	// faz requisicao GET para a OpenWeather
	resp, err := client.Get(apiURL)
	if err != nil {
		appLogger.Error("Erro ao chamar OpenWeather Forecast: " + err.Error())

		return model.ForecastToolResponse{}, &ForecastError{
			StatusCode:  http.StatusBadGateway, // codigo 502
			UserMessage: "Erro ao consultar serviço externo de previsão.",
			LogMessage:  "Erro ao chamar OpenWeather Forecast: " + err.Error(),
		}
	}

	// garante fechamento do corpo da resposta ao final da funcao
	defer resp.Body.Close()

	// tatamento de codigos de erro
	if resp.StatusCode == http.StatusNotFound {
		appLogger.Warning("Cidade não encontrada na OpenWeather Forecast: " + city)

		return model.ForecastToolResponse{}, &ForecastError{
			StatusCode:  http.StatusNotFound, // codigo 404
			UserMessage: "Cidade não encontrada.",
			LogMessage:  "Cidade não encontrada na OpenWeather Forecast: " + city,
		}
	}

	if resp.StatusCode != http.StatusOK {
		appLogger.Error(fmt.Sprintf("OpenWeather retornou status inesperado em /forecast: %d", resp.StatusCode))

		return model.ForecastToolResponse{}, &ForecastError{
			StatusCode:  http.StatusBadGateway, // codigo 502
			UserMessage: "Erro retornado pela OpenWeather.",
			LogMessage:  fmt.Sprintf("OpenWeather retornou status inesperado em /forecast: %d", resp.StatusCode),
		}
	}

	// variavel que recebera JSON bruto da OpenWeather
	var forecast model.OpenWeatherForecastResponse

	// faz parse do JSON bruto diretamente para a struct OpenWeatherForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		appLogger.Error("Erro ao decodificar JSON da OpenWeather Forecast: " + err.Error())

		return model.ForecastToolResponse{}, &ForecastError{
			StatusCode:  http.StatusInternalServerError,
			UserMessage: "Erro ao processar resposta da OpenWeather.",
			LogMessage:  "Erro ao decodificar JSON da OpenWeather Forecast: " + err.Error(),
		}
	}

	// transforma o JSON bruto da OpenWeather em JSON limpo para o agent/frontend
	response := mapper.ToForecastToolResponse(forecast, days)
	appLogger.Info(fmt.Sprintf("Forecast processado com sucesso para cidade: %s, itens retornados: %d", response.City, len(response.Items)))

	return response, nil
}