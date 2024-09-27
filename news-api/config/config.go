package config

import (
	"log"
	"os"
)

func LoadConfig() {
	log.Println("Loading configuration...")
	// Здесь можно загружать конфигурацию из переменных окружения или файлов
	if _, exists := os.LookupEnv("PORT"); !exists {
		os.Setenv("PORT", "8080")
	}
}
