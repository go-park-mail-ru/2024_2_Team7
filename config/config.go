package config

import (
	"log"
	"os"
)

const (
	DefaultPort = "5500"
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
