//responde /forecast

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url" //faz parsing de urls e implementa escape de query strings
	"strconv"
	"time"
	appLogger "web-service-gin/src/infrastructure/logger"
	"web-service-gin/src/config"
	"web-service-gin/src/infrastructure/model"

	"github.com/gin-gonic/gin"
)

//controller

// Cria um handler resoponsavel pelo endpoint -> GET /forecast?city={cidade}&days={n}
// Recebe config pronta (variaveis de ambiente)

func GetForecast(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		appLogger.Info("Requisição recebida em /forecast")

		city := c.Query("city") //le parametro city da URL
		daysParam := c.DefaultQuery("days", "3") //le parametro days da URL (valor padrao de 3 dias)

		if city == "" {
			appLogger.Warning("Parâmetro city não informado em /forecast")
			c.JSON(http.StatusBadRequest, gin.H{ // codigo 400
				"error": "O parâmetro 'city' é obrigatório.",
			})
			return
		}

		// service 
		// converte os dias de string para int 
		days, err := strconv.Atoi(daysParam)
		if err != nil || days < 1 || days > 5 {
			appLogger.Warning("Parâmetro days inválido em /forecast: " + daysParam)
			c.JSON(http.StatusBadRequest, gin.H{ // codigo 400
				"error": "O parâmetro 'days' deve ser um número entre 1 e 5.",
			})
			return
		}

		// service
		//verificar se o OpenWeather foi carregada corretamente
		if cfg.OpenWeatherAPIKey == "" {
			appLogger.Critical("OPENWEATHER_API_KEY não configurada")
			c.JSON(http.StatusInternalServerError, gin.H{ // codigo 500
				"error": "OPENWEATHER_API_KEY não configurada.",
			})
			return
		}

		//verifica se a URL do endpoint forecast foi configurada
		if cfg.OpenWeatherForecastURL == "" {
			appLogger.Critical("OPENWEATHER_FORECAST_URL não configurada")
			c.JSON(http.StatusInternalServerError, gin.H{ // codigo 500
				"error": "OPENWEATHER_FORECAST_URL não configurada.",
			})
			return
		}

		//trata caracteres especiais para montar a URL da API externa corretamente
		encodedCity := url.QueryEscape(city)


		//service -> puxar de formatter ( classe weatherFormatter -> func forecast_formatter) 
		//monta a URL final do OpenWeather
		apiURL := fmt.Sprintf(
			"%s?q=%s&appid=%s&units=metric&lang=pt_br",
			cfg.OpenWeatherForecastURL,
			encodedCity,
			cfg.OpenWeatherAPIKey,
		)

		//service
		//cria cliente http com timeout
		client := http.Client{
			Timeout: 10 * time.Second,
		}

		appLogger.Info(fmt.Sprintf("Chamando OpenWeather Forecast para cidade: %s, dias: %d", city, days))


		//faz requisicao GET para a OpenWeather
		resp, err := client.Get(apiURL)

		if err != nil {
			appLogger.Error("Erro ao chamar OpenWeather Forecast: " + err.Error())
			c.JSON(http.StatusBadGateway, gin.H{ // codigo 502
				"error": "Erro ao consultar serviço externo de previsão.",
			})
			return
		}

		//garante fechamento do corpo da resposta ao final da funcao
		defer resp.Body.Close()

		//service
		//tatamento de codigos de erro
		if resp.StatusCode == http.StatusNotFound {
			appLogger.Warning("Cidade não encontrada na OpenWeather Forecast: " + city)
			c.JSON(http.StatusNotFound, gin.H{ // codigo 404
				"error": "Cidade não encontrada.",
			})
			return
		}

		// service (puxar erro p usuario no controller)
		if resp.StatusCode != http.StatusOK {
			appLogger.Error(fmt.Sprintf("OpenWeather retornou status inesperado em /forecast: %d", resp.StatusCode))
			c.JSON(http.StatusBadGateway, gin.H{ // codigo 502
				"error": "Erro retornado pela OpenWeather.",
			})
			return
		}

		//service
		//variavel que recebera JSON bruto da OpenWeather
		var forecast model.OpenWeatherForecastResponse //

		//faz parse do JSON bruto diretamente para a struct OpenWeatherForecastResponse
		if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
			appLogger.Error("Erro ao decodificar JSON da OpenWeather Forecast: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao processar resposta da OpenWeather.",
			})
			return
		}

		//calcula quantidade de registros pode dia (8)
		limit := days * 8
		if limit > len(forecast.List) {
			limit = len(forecast.List)
		}

		//cria slice vazio de FOrecastItem com capacidade inicial igual ao limite (evitar realocacoes desnecessarias)
		items := make([]model.ForecastItem, 0, limit)

		//percorre apenas pelos registros necessarios, de acordo com numero de dias
		for i := 0; i < limit; i++ {
			item := forecast.List[i]

			description := ""
			if len(item.Weather) > 0 {
				description = item.Weather[0].Description
			}

			//transforma o item bruto da OpenWeather em uma estrutura mais simples para o agent e frontend consumirem
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

		//monta resposta final limpa -> resposta que a tool do agent vai receber
		response := model.ForecastToolResponse{
			Type:          "weather_forecast",
			City:          forecast.City.Name,
			Country:       forecast.City.Country,
			DaysRequested: days,
			Items:         items,
		}

		// responseJSON, err := json.MarshalIndent(response, "", " ")
		// 	if err != nil {
		// 		appLogger.Error("Erro ao converter resposta forecast para JSON: " + err.Error())
		// 	} else {
		// 		appLogger.Info("JSON limpo retornado para a tool: " + string(responseJSON))
		// }
		
		//retorna a resposta do JSON com status 200 (OK)
		c.JSON(http.StatusOK, response)
		appLogger.Info(fmt.Sprintf("Forecast processado com sucesso para cidade: %s, itens retornados: %d", response.City, len(response.Items)))
	}
}