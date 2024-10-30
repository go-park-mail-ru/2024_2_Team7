package sessionRepository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"kudago/internal/models"
)

const (
	ExpirationTime = 24 * time.Hour
)

type SessionDB struct {
	client *redis.Client
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	URL      string
	DB       int
	PoolSize int
}

func NewDB(config *RedisConfig) *SessionDB {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	return &SessionDB{
		client: redisClient,
	}
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

func (db *SessionDB) CreateSession(ctx context.Context, ID int) (models.Session, error) {
	sessionToken := generateSessionToken()
	expiration := time.Now().Add(ExpirationTime)

	session := models.Session{
		UserID:  ID,
		Token:   sessionToken,
		Expires: expiration,
	}

	err := db.client.Set(ctx, sessionToken, session.UserID, ExpirationTime).Err()
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

func (db *SessionDB) CheckSession(ctx context.Context, cookie string) (models.Session, error) {
	ID, err := db.client.Get(ctx, cookie).Result()
	if err == redis.Nil {
		return models.Session{}, models.ErrUserNotFound
	}

	if err != nil {
		return models.Session{}, err
	}

	userID, err := strconv.Atoi(ID)
	if err != nil {
		return models.Session{}, err
	}

	session := models.Session{
		UserID:  userID,
		Token:   cookie,
		Expires: time.Now().Add(ExpirationTime),
	}

	return session, nil
}

func (db *SessionDB) DeleteSession(ctx context.Context, token string) error {
	err := db.client.Del(ctx, token).Err()
	if err != nil {
		return err
	}
	return nil
}

func generateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
