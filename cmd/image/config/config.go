package config

import (
	"errors"
	"os"

	imageRepository "kudago/internal/image/repository"

	"github.com/joho/godotenv"
)

type Config struct {
	ImageConfig imageRepository.ImageConfig
	ServiceAddr string
}

func LoadConfig() (Config, error) {
	var conf Config
	err := godotenv.Load()
	if err != nil {
		return conf, err
	}

	conf.ImageConfig = imageRepository.ImageConfig{
		Path: "./static/images",
	}

	conf.ServiceAddr = os.Getenv("IMAGE_SERVICE_ADDR")
	if conf.ServiceAddr == "" {
		return Config{}, errors.New("Failed to get service address")
	}

	return conf, nil
}
