package config

import (
	"errors"
	"os"

	"kudago/internal/repository/postgres"

	sessionRepository "kudago/internal/repository/redis"
)

type Config struct {
	PostgresConfig postgres.PostgresConfig
	RedisConfig    sessionRepository.RedisConfig
	ServiceAddr    string
}

func LoadConfig() (Config, error) {
	var conf Config

	redisConfig, err := sessionRepository.GetRedisConfig()
	if err != nil {
		return Config{}, errors.New("Failed to connect to the redis database")
	}
	conf.RedisConfig = *redisConfig

	postgresConfig, err := postgres.GetPostgresConfig()
	if err != nil {
		return Config{}, errors.New("Failed to connect to the postgres database")
	}
	conf.PostgresConfig = postgresConfig

	conf.ServiceAddr = os.Getenv("AUTH_SERVICE_ADDR")
	if conf.ServiceAddr == "" {
		return Config{}, errors.New("Failed to get service address")
	}

	return conf, nil
}
