// Carrega variaveis de ambiente

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenWeatherAPIKey      string
	OpenWeatherCurrentURL  string
	OpenWeatherForecastURL string
	ServerAddress          string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado. Usando variáveis do sistema.")
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost:8080"
	}

	return Config{
		OpenWeatherAPIKey:      os.Getenv("OPENWEATHER_API_KEY"),
		OpenWeatherCurrentURL:  os.Getenv("OPENWEATHER_CURRENT_URL"),
		OpenWeatherForecastURL: os.Getenv("OPENWEATHER_FORECAST_URL"),
		ServerAddress:          serverAddress,
	}
}