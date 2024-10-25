package db

import (
	"context"
	"fmt"
	"os"

	"kudago/internal/http/utils"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitDB(logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error in loading .env: %v", err)
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	dbConf, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse db URL: %v", err)
	}
	dbConf.ConnConfig.Logger = zapLogger{logger}
	pool, err := pgxpool.ConnectConfig(context.Background(), dbConf)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	return pool, nil
}

func InitRedis() (*redis.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error in loading .env: %v", err)
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisURL := fmt.Sprintf("%s:%s", redisHost, redisPort)

	return redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       0,
	}), nil
}

type zapLogger struct {
	logger *zap.SugaredLogger
}

func (z zapLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	requestID, ok := utils.GetRequestIDFromContext(ctx)
	if ok {
		if len(data) > 0 {
			z.logger.Infow(msg,
				"request_id", requestID,
				"level", level,
				"data", data,
			)
		} else {
			z.logger.Infow(msg,
				"request_id", requestID,
				"level", level,
			)
		}
	}
}
