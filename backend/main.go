package main

import (
	"log"

	"web-service-gin/internal/config"
	"web-service-gin/internal/handlers"

	"github.com/gin-gonic/gin"
)

// 1. Carregar configuração
// 2. Criar o router
// 3. Registrar as rotas
// 4. Subir o servidor

func main() {
	cfg := config.Load()

	router := gin.Default()

	router.GET("/health", handlers.HealthCheck)
	router.GET("/weather", handlers.GetWeather(cfg))
	router.GET("/forecast", handlers.GetForecast(cfg))

	log.Printf("Servidor rodando em http://%s\n", cfg.ServerAddress)

	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}