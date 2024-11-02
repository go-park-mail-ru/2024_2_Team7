package sessionRepository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"kudago/internal/models"
	redisDB "kudago/internal/repository/redis"
)

const (
	expirationTime = 24 * time.Hour
)

type SessionDB struct {
	client *redis.Client
}

func NewDB(config *redisDB.RedisConfig) *SessionDB {
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

func (db *SessionDB) CreateSession(ctx context.Context, ID int) (models.Session, error) {
	sessionToken := generateSessionToken()
	expiration := time.Now().Add(expirationTime)

	session := models.Session{
		UserID:  ID,
		Token:   sessionToken,
		Expires: expiration,
	}

	err := db.client.Set(ctx, sessionToken, session.UserID, expirationTime).Err()
	if err != nil {
		return models.Session{}, errors.Wrap(err, models.LevelDB)
	}
	return session, nil
}

func (db *SessionDB) CheckSession(ctx context.Context, cookie string) (models.Session, error) {
	ID, err := db.client.Get(ctx, cookie).Result()
	if err == redis.Nil {
		return models.Session{}, models.ErrUserNotFound
	}

	if err != nil {
		return models.Session{}, errors.Wrap(err, models.LevelDB)
	}

	userID, err := strconv.Atoi(ID)
	if err != nil {
		return models.Session{}, errors.Wrap(err, models.LevelDB)
	}

	session := models.Session{
		UserID:  userID,
		Token:   cookie,
		Expires: time.Now().Add(expirationTime),
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
