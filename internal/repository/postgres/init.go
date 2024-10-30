package postgres

import (
	"context"
	"fmt"
	"os"

	"kudago/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventsConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
	URL      string
}

func GetPostgresConfig() (EventsConfig, error) {
	var config EventsConfig
	config.User = os.Getenv("POSTGRES_USER")
	config.Password = os.Getenv("POSTGRES_PASSWORD")
	config.Host = os.Getenv("POSTGRES_HOST")
	config.Port = os.Getenv("POSTGRES_PORT")
	config.DB = os.Getenv("POSTGRES_DB")
	config.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.User, config.Password, config.Host, config.Port, config.DB)

	return config, nil
}

func InitPostgres(config EventsConfig, logger *logger.Logger) (*pgxpool.Pool, error) {
	dbConf, err := pgxpool.ParseConfig(config.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse db URL: %v", err)
	}
	dbConf.ConnConfig.Tracer = logger

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConf)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	postgresPing := pool.Ping(context.Background())
	if err != nil || postgresPing != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}
	return pool, nil
}
