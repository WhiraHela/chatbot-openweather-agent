// Responde /health

package controller

import (
	"net/http" //pacote que fornece implemetacoes cliente e servidor HTTP
	appLogger "web-service-gin/src/infrastructure/logger"
	"github.com/gin-gonic/gin" //framework HTTP para Go (aplicacoes web, APIs REST, microservicos)
)

//Corresponde ao endpoint /health -> GET /health

// Verifica se a API esta rodando
func HealthCheck(c *gin.Context) {
	appLogger.Info("Health check recebido")
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}