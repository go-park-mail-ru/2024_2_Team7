package config

import (
	"errors"
	"fmt"
	"os"

	"kudago/internal/repository/postgres"

	imageRepository "kudago/internal/repository/images"
	sessionRepository "kudago/internal/repository/redis"
)

const (
	DefaultPort = "8080"
)

type Config struct {
	Port           string
	PostgresConfig postgres.PostgresConfig
	RedisConfig    sessionRepository.RedisConfig
	ImageConfig    imageRepository.ImageConfig
}

func LoadConfig() (Config, error) {
	var conf Config

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	conf.Port = port

	redisConfig, err := sessionRepository.GetRedisConfig()
	if err != nil {
		return Config{}, errors.New("Failed to connect to the redis database")
	}
	conf.RedisConfig = *redisConfig

	postgresConfig, err := postgres.GetPostgresConfig()
	if err != nil {
		fmt.Println(err)
		return Config{}, errors.New("Failed to connect to the postgres database")
	}
	conf.PostgresConfig = postgresConfig

	conf.ImageConfig = imageRepository.ImageConfig{
		Path: "./static/images",
	}

	return conf, nil
}
