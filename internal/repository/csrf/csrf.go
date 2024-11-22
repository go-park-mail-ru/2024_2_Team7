package csrf

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/redis/go-redis/v9"
	"kudago/internal/models"
)

const CSRFTokenExpTime = 15 * time.Minute

type csrfDB struct {
	client *redis.Client
}

func NewDB(client *redis.Client) *csrfDB {
	return &csrfDB{client: client}
}

func (db *csrfDB) CreateCSRF(ctx context.Context, encryptionKey []byte, s *models.Session) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	tokenExpTime := time.Now().Add(CSRFTokenExpTime)
	td := models.TokenData{
		SessionToken: s.Token,
		UserID:       s.UserID,
		Exp:          tokenExpTime,
	}
	data, _ := json.Marshal(td)
	ciphertext := aesgcm.Seal(nil, nonce, data, nil)

	res := append(nonce, ciphertext...)
	token := base64.StdEncoding.EncodeToString(res)

	err = db.client.Set(ctx, prefixedKey(s.Token), token, CSRFTokenExpTime).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store token in Redis: %v", err)
	}

	return token, nil
}

func (db *csrfDB) GetCSRF(ctx context.Context, s *models.Session) (string, error) {
	storedToken, err := db.client.Get(ctx, prefixedKey(s.Token)).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("token not found in Redis")
	} else if err != nil {
		return "", fmt.Errorf("failed to get token from Redis: %v", err)
	}
	return storedToken, nil
}

func prefixedKey(key string) string {
	return fmt.Sprintf("%s:%s", key, "csrf")
}
