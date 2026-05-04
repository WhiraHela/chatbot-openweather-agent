//responde /weather

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	appLogger "web-service-gin/src/infrastructure/logger"

	"web-service-gin/src/config"
	"web-service-gin/src/infrastructure/model"

	"github.com/gin-gonic/gin"
)

//Reponsavel pelo endpoint /weather -> GET /weather?city={cidade}

//Busca o clima atual de uma cidade:
//	Chama OpenWeather, recebe JSON bruto, transforma em JSON mais limpo, devolve para a tool get_current_weather no frontend
func GetWeather(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		//le parametro city da URL
		city := c.Query("city")

		appLogger.Info("Requisição recebida em /weather")
		if city == "" { //verifica se o usuario enviou uma cidade
			appLogger.Warning("Parâmetro city não informado em /weather")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "O parâmetro 'city' é obrigatório.",
			})
			return
		}

		if cfg.OpenWeatherAPIKey == "" { //verifica se a chave da OpenWeather foi carregada
			appLogger.Critical("OPENWEATHER_API_KEY não configurada")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "OPENWEATHER_API_KEY não configurada.",
			})
			return
		}

		if cfg.OpenWeatherCurrentURL == "" { //verifica se a URL do endpoint de clima atual foi carregada
			appLogger.Critical("OPENWEATHER_CURRENT_URL não configurada")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "OPENWEATHER_CURRENT_URL não configurada.",
			})
			return
		}

		appLogger.Info("Chamando OpenWeather Current para cidade: " + city)
		//codifica cidade para uso seguro dentro da URL
		encodedCity := url.QueryEscape(city)

		apiURL := fmt.Sprintf(
			"%s?q=%s&appid=%s&units=metric&lang=pt_br",
			cfg.OpenWeatherCurrentURL,
			encodedCity,
			cfg.OpenWeatherAPIKey,
		)

		client := http.Client{ //client http com timeout
			Timeout: 10 * time.Second,
		}

		//faz chamada para a OpenWeather
		resp, err := client.Get(apiURL)
		if err != nil {
			appLogger.Error("Erro ao chamar OpenWeather Current: " + err.Error())
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Erro ao consultar serviço externo de clima.",
			})
			return
		}
		//fecha corpo da resposta ao final da funcao
		defer resp.Body.Close()

		// tratamento de cidade nao encontrada
		if resp.StatusCode == http.StatusNotFound {
			appLogger.Warning("Cidade não encontrada na OpenWeather: " + city)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cidade não encontrada.",
			})
			return
		}

		//tratamento de erro geral
		if resp.StatusCode != http.StatusOK {
			appLogger.Error(fmt.Sprintf("OpenWeather retornou status inesperado em /weather: %d", resp.StatusCode))
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Erro retaornado pela OpenWeather.",
			})
			return
		}

		//struct que recebe JSON bruto da OpenWeather
		var openWeather model.OpenWeatherResponse

		// decodifica o JSON da resposta para a struct Go
		if err := json.NewDecoder(resp.Body).Decode(&openWeather); err != nil {
			appLogger.Error("Erro ao decodificar JSON da OpenWeather Current: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao processar resposta da OpenWeather.",
			})
			return
		}

		//retorna weather como array - so acessamos o primeiro item se ele existir
		description := ""
		if len(openWeather.Weather) > 0 {
			description = openWeather.Weather[0].Description
		}

		//monta resposta limpa para o agent
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

		//retorna resposta final em JSON
		c.JSON(http.StatusOK, response)
		appLogger.Info("Resposta de clima atual processada com sucesso para cidade: " + response.City)
	}
}