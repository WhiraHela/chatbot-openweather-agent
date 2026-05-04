// responde /forecast

package controller

import (
	"errors"
	"net/http"
	"strconv"

	"web-service-gin/src/config"
	appLogger "web-service-gin/src/infrastructure/logger"
	"web-service-gin/src/infrastructure/service"

	"github.com/gin-gonic/gin"
)

// controller

// Cria um handler resoponsavel pelo endpoint -> GET /forecast?city={cidade}&days={n}
// Recebe config pronta (variaveis de ambiente)
func GetForecast(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		appLogger.Info("Requisição recebida em /forecast")

		city := c.Query("city")                  // le parametro city da URL
		daysParam := c.DefaultQuery("days", "3") // le parametro days da URL (valor padrao de 3 dias)

		if city == "" {
			appLogger.Warning("Parâmetro city não informado em /forecast")
			c.JSON(http.StatusBadRequest, gin.H{ // codigo 400
				"error": "O parâmetro 'city' é obrigatório.",
			})
			return
		}

		// converte os dias de string para int
		// essa conversão fica no controller porque days vem da URL como texto
		days, err := strconv.Atoi(daysParam)
		if err != nil || days < 1 || days > 5 {
			appLogger.Warning("Parâmetro days inválido em /forecast: " + daysParam)
			c.JSON(http.StatusBadRequest, gin.H{ // codigo 400
				"error": "O parâmetro 'days' deve ser um número entre 1 e 5.",
			})
			return
		}

		// chama a camada de service, que concentra a regra de negócio da previsão
		response, err := service.GetForecast(cfg, city, days)
		if err != nil {
			var forecastErr *service.ForecastError

			if errors.As(err, &forecastErr) {
				c.JSON(forecastErr.StatusCode, gin.H{
					"error": forecastErr.UserMessage,
				})
				return
			}

			appLogger.Error("Erro inesperado em /forecast: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro interno ao processar previsão.",
			})
			return
		}

		// retorna a resposta do JSON com status 200 (OK)
		c.JSON(http.StatusOK, response)
		appLogger.Info("Forecast processado com sucesso para cidade: " + response.City)
	}
}