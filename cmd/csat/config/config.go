package config

import (
	"errors"
	"os"

	"kudago/internal/repository/postgres"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresConfig postgres.PostgresConfig
	ServiceAddr    string
}

func LoadConfig() (Config, error) {
	var conf Config
	err := godotenv.Load()
	if err != nil {
		return conf, err
	}

	postgresConfig, err := postgres.GetPostgresConfig()
	if err != nil {
		return Config{}, errors.New("Failed to connect to the postgres database")
	}
	conf.PostgresConfig = postgresConfig

	conf.ServiceAddr = os.Getenv("CSAT_SERVICE_ADDR")
	if conf.ServiceAddr == "" {
		return Config{}, errors.New("Failed to get service address")
	}

	return conf, nil
}
