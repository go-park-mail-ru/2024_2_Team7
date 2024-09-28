package config

import (
	"log"
	"os"
)

func LoadConfig() {
	log.Println("Loading configuration...")
	if _, exists := os.LookupEnv("PORT"); !exists {
		os.Setenv("PORT", "8080")
	}
}
