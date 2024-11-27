package redisDB

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	URL      string
	DB       int
	PoolSize int
}

func GetRedisConfig() (*RedisConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error in loading .env: %v", err)
	}

	var config RedisConfig
	config.Password = os.Getenv("REDIS_PASSWORD")
	config.Host = os.Getenv("REDIS_HOST")
	config.Port = os.Getenv("REDIS_PORT")
	config.URL = fmt.Sprintf("%s:%s", config.Host, config.Port)

	return &config, nil
}
