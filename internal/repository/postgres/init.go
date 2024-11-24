package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"

	"kudago/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
	URL      string
}

func GetPostgresConfig() (PostgresConfig, error) {
	var config PostgresConfig
	config.User = os.Getenv("POSTGRES_USER")
	config.Password = os.Getenv("POSTGRES_PASSWORD")
	config.Host = os.Getenv("POSTGRES_HOST")
	config.Port = os.Getenv("POSTGRES_PORT")
	config.DB = os.Getenv("POSTGRES_DB")
	config.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.User, config.Password, config.Host, config.Port, config.DB)

	return config, nil
}

func InitPostgres(config PostgresConfig, logger *logger.Logger) (*pgxpool.Pool, error) {
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
	if postgresPing != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	if err := RunMigrations(config.URL); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return pool, nil
}

func RunMigrations(dbURL string) error {
	migrationsDir := os.Getenv("MIGRATION_FOLDER")
	if migrationsDir == "" {
		return fmt.Errorf("MIGRATION_FOLDER environment variable is not set")
	}

	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("unable to open db: %v", err)
	}
	defer sqlDB.Close()

	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}
