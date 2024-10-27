package sessionRepository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"kudago/internal/models"
)

const (
	ExpirationTime = 24 * time.Hour
)

type SessionDB struct {
	client *redis.Client
}

func NewDB(client *redis.Client) *SessionDB {
	return &SessionDB{
		client: client,
	}
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

func (db *SessionDB) CheckSession(ctx context.Context, cookie string) (models.Session, bool) {
	ID, err := db.client.Get(ctx, cookie).Result()
	if err == redis.Nil {
		return models.Session{}, false
	}

	if err != nil {
		return models.Session{}, false
	}

	userID, err := strconv.Atoi(ID)
	if err != nil {
		return models.Session{}, false
	}

	session := models.Session{
		UserID:  userID,
		Token:   cookie,
		Expires: time.Now().Add(ExpirationTime),
	}

	return session, true
}

func (db *SessionDB) DeleteSession(ctx context.Context, token string) {
	db.client.Del(ctx, token)
}

func generateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
