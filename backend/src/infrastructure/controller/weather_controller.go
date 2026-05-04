// responde /weather

package controller

import (
	"errors"
	"net/http"

	"web-service-gin/src/config"
	appLogger "web-service-gin/src/infrastructure/logger"
	"web-service-gin/src/infrastructure/service"

	"github.com/gin-gonic/gin"
)

// Reponsavel pelo endpoint /weather -> GET /weather?city={cidade}

// Busca o clima atual de uma cidade:
// Chama a camada de service, que consulta a OpenWeather,
// recebe JSON bruto, transforma em JSON mais limpo,
// e devolve para a tool get_current_weather no frontend.
func GetWeather(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		appLogger.Info("Requisição recebida em /weather")

		// le parametro city da URL
		city := c.Query("city")

		if city == "" { // verifica se o usuario enviou uma cidade
			appLogger.Warning("Parâmetro city não informado em /weather")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "O parâmetro 'city' é obrigatório.",
			})
			return
		}

		// chama a camada de service, que concentra a regra de negócio do clima atual
		response, err := service.GetCurrentWeather(cfg, city)
		if err != nil {
			var weatherErr *service.WeatherError

			if errors.As(err, &weatherErr) {
				c.JSON(weatherErr.StatusCode, gin.H{
					"error": weatherErr.UserMessage,
				})
				return
			}

			appLogger.Error("Erro inesperado em /weather: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro interno ao processar clima atual.",
			})
			return
		}

		// retorna resposta final em JSON
		c.JSON(http.StatusOK, response)
		appLogger.Info("Resposta de clima atual processada com sucesso para cidade: " + response.City)
	}
}