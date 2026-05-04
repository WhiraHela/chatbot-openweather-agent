// Carrega variaveis de ambiente

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config centraliza todas as configurações da aplicação.
type Config struct {
	OpenWeatherAPIKey      string
	OpenWeatherCurrentURL  string
	OpenWeatherForecastURL string
	ServerAddress          string
}

// Load carrega as variáveis de ambiente do arquivo .env
// e devolve uma struct Config preenchida.
func Load() Config {
	if err := godotenv.Load(); err != nil { // se nao tiver erro para carregar o dotenv
		log.Println("Arquivo .env não encontrado. Usando variáveis do sistema.")
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost:8080"
	}

	// Retorna todas as configurações necessárias para a aplicação.
	// Os handlers vão receber essa struct e usar esses valores para chamar a OpenWeather.
	return Config{
		OpenWeatherAPIKey:      os.Getenv("OPENWEATHER_API_KEY"),
		OpenWeatherCurrentURL:  os.Getenv("OPENWEATHER_CURRENT_URL"),
		OpenWeatherForecastURL: os.Getenv("OPENWEATHER_FORECAST_URL"),
		ServerAddress:          serverAddress,
	}
}