package main

import (
	appLogger "web-service-gin/src/logger"

	"web-service-gin/src/config"
	"web-service-gin/src/controller"

	"github.com/gin-gonic/gin"
)

// 1. Carregar configuração
// 2. Criar o router
// 3. Registrar as rotas
// 4. Subir o servidor

func main() {
	cfg := config.Load()

	appLogger.Info("Configurações carregadas")
	appLogger.Info("Inicializando servidor Gin")

	router := gin.Default()

	//registrando rotas
	router.GET("/health", controller.HealthCheck)
	router.GET("/weather", controller.GetWeather(cfg))
	router.GET("/forecast", controller.GetForecast(cfg))

	appLogger.Info("Rotas registradas: /health, /weather, /forecast")
	appLogger.Info("Servidor rodando em http://" + cfg.ServerAddress)

	// trata erro de execucao do servidor
	if err := router.Run(cfg.ServerAddress); err != nil {
		appLogger.Critical("Erro crítico ao iniciar servidor: " + err.Error())
	}
}





