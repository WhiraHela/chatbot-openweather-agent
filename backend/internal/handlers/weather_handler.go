//responde /weather

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"web-service-gin/internal/config"
	"web-service-gin/internal/models"

	"github.com/gin-gonic/gin"
)

func GetWeather(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		city := c.Query("city")

		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "O parâmetro 'city' é obrigatório.",
			})
			return
		}

		if cfg.OpenWeatherAPIKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "OPENWEATHER_API_KEY não configurada.",
			})
			return
		}

		if cfg.OpenWeatherCurrentURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "OPENWEATHER_CURRENT_URL não configurada.",
			})
			return
		}

		encodedCity := url.QueryEscape(city)

		apiURL := fmt.Sprintf(
			"%s?q=%s&appid=%s&units=metric&lang=pt_br",
			cfg.OpenWeatherCurrentURL,
			encodedCity,
			cfg.OpenWeatherAPIKey,
		)

		client := http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(apiURL)
		if err != nil {
			log.Println("Erro ao chamar OpenWeather:", err)

			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Erro ao consultar serviço externo de clima.",
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cidade não encontrada.",
			})
			return
		}

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Erro retornado pela OpenWeather.",
			})
			return
		}

		var openWeather models.OpenWeatherResponse

		if err := json.NewDecoder(resp.Body).Decode(&openWeather); err != nil {
			log.Println("Erro ao decodificar JSON:", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao processar resposta da OpenWeather.",
			})
			return
		}

		description := ""
		if len(openWeather.Weather) > 0 {
			description = openWeather.Weather[0].Description
		}

		response := models.WeatherToolResponse{
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

		c.JSON(http.StatusOK, response)
	}
}