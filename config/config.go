package config

import (
	"log"
	"os"
)

const (
	DefaultPort = "8080"
)

func LoadConfig() string {
	log.Println("Loading configuration...")
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	log.Printf("Используется порт: %s", port)
	return port
}

func LoadEncriptionKey() []byte {
	key := os.Getenv("ENCRYPTION_KEY")
	return []byte(key)
}
