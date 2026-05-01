//responde /forecast

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"web-service-gin/internal/config"
	"web-service-gin/internal/models"

	"github.com/gin-gonic/gin"
)

func GetForecast(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		city := c.Query("city")
		daysParam := c.DefaultQuery("days", "3")

		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "O parâmetro 'city' é obrigatório.",
			})
			return
		}

		days, err := strconv.Atoi(daysParam)
		if err != nil || days < 1 || days > 5 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "O parâmetro 'days' deve ser um número entre 1 e 5.",
			})
			return
		}

		if cfg.OpenWeatherAPIKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "OPENWEATHER_API_KEY não configurada.",
			})
			return
		}

		if cfg.OpenWeatherForecastURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "OPENWEATHER_FORECAST_URL não configurada.",
			})
			return
		}

		encodedCity := url.QueryEscape(city)

		apiURL := fmt.Sprintf(
			"%s?q=%s&appid=%s&units=metric&lang=pt_br",
			cfg.OpenWeatherForecastURL,
			encodedCity,
			cfg.OpenWeatherAPIKey,
		)

		client := http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(apiURL)
		if err != nil {
			log.Println("Erro ao chamar OpenWeather forecast:", err)

			c.JSON(http.StatusBadGateway, gin.H{
				"error": "Erro ao consultar serviço externo de previsão.",
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

		var forecast models.OpenWeatherForecastResponse

		if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
			log.Println("Erro ao decodificar JSON do forecast:", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao processar resposta da OpenWeather.",
			})
			return
		}

		limit := days * 8
		if limit > len(forecast.List) {
			limit = len(forecast.List)
		}

		items := make([]models.ForecastItem, 0, limit)

		for i := 0; i < limit; i++ {
			item := forecast.List[i]

			description := ""
			if len(item.Weather) > 0 {
				description = item.Weather[0].Description
			}

			forecastItem := models.ForecastItem{
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

		response := models.ForecastToolResponse{
			Type:          "weather_forecast",
			City:          forecast.City.Name,
			Country:       forecast.City.Country,
			DaysRequested: days,
			Items:         items,
		}

		c.JSON(http.StatusOK, response)
	}
}