package db

import (
	"context"
	"fmt"
	"os"

	"kudago/internal/http/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	dbConf.ConnConfig.Tracer = &zapLogger{logger}
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConf)
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       0,
		PoolSize: 10,
	})

	redisPing := redisClient.Ping(context.Background())
	if redisPing.Err() != nil {
		return nil, redisPing.Err()
	}
	return redisClient, nil
}

type zapLogger struct {
	logger *zap.SugaredLogger
}

func (z *zapLogger) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	requestID := utils.GetRequestIDFromContext(ctx)
	z.logger.Infow("Query",
		"request_id", requestID,
		"sql", data.SQL,
		"args", data.Args,
	)
	return ctx
}

func (z *zapLogger) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	requestID := utils.GetRequestIDFromContext(ctx)
	if data.Err != nil {
		z.logger.Errorw("Query failed",
			"request_id", requestID,
			"commandTag", data.CommandTag,
			"args", data.Err,
		)
	}
}
