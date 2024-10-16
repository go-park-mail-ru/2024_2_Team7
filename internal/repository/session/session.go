package sessionRepository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"kudago/internal/models"
)

const (
	ExpirationTime = 24 * time.Hour
)

type SessionDB struct {
	mu       *sync.RWMutex
	sessions map[string]models.Session
}

func NewDB() *SessionDB {
	return &SessionDB{
		sessions: make(map[string]models.Session),
		mu:       &sync.RWMutex{},
	}
}

func (db *SessionDB) CreateSession(ctx context.Context, ID int) *models.Session {
	sessionToken := generateSessionToken()
	expiration := time.Now().Add(ExpirationTime)

	session := models.Session{
		UserID:  ID,
		Token:   sessionToken,
		Expires: expiration,
	}
	db.mu.Lock()
	db.sessions[sessionToken] = session
	db.mu.Unlock()
	return &session
}

func (db SessionDB) CheckSession(ctx context.Context, cookie string) (*models.Session, bool) {
	db.mu.RLock()
	session, exists := db.sessions[cookie]
	db.mu.RUnlock()
	if !exists || session.Expires.Before(time.Now()) {
		return nil, false
	}
	return &session, true
}

func (db *SessionDB) DeleteSession(ctx context.Context, token string) {
	db.mu.Lock()
	delete(db.sessions, token)
	db.mu.Unlock()
}

func generateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
