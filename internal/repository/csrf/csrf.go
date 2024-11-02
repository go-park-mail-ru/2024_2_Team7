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

	err = db.client.Set(ctx, PrefixedKey(s.Token), token, CSRFTokenExpTime).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store token in Redis: %v", err)
	}

	return token, nil
}

func (db *csrfDB) CheckCSRF(ctx context.Context, encryptionKey []byte, s *models.Session, inputToken string) (bool, error) {
	storedToken, err := db.client.Get(ctx, PrefixedKey(s.Token)).Result()
	if err == redis.Nil {
		return false, fmt.Errorf("token not found in Redis")
	} else if err != nil {
		return false, fmt.Errorf("failed to get token from Redis: %v", err)
	}

	if storedToken != inputToken {
		return false, fmt.Errorf("invalid token")
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return false, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return false, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(inputToken)
	if err != nil {
		return false, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return false, fmt.Errorf("short ciphertext")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, fmt.Errorf("decrypt fail: %v", err)
	}

	td := &models.TokenData{}
	err = json.Unmarshal(plaintext, &td)
	if err != nil {
		return false, fmt.Errorf("bad json: %v", err)
	}

	if td.Exp.Unix() < time.Now().Unix() {
		return false, fmt.Errorf("token expired")
	}

	return s.Token == td.SessionToken && s.UserID == td.UserID, nil
}

func PrefixedKey(key string) string {
	return fmt.Sprintf("%s:%s", key, "csrf")
}
