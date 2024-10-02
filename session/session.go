package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

const (
	SessionToken   = "session_token"
	ExpirationTime = 24 * time.Hour
)

type Session struct {
	Username string
	Token    string
	Expires  time.Time
}

type SessionDB struct {
	mu       *sync.RWMutex
	sessions map[string]Session
}

func NewSessionDB() *SessionDB {
	return &SessionDB{
		sessions: make(map[string]Session),
		mu:       &sync.RWMutex{},
	}
}

func (db *SessionDB) CreateSession(username string) *Session {
	sessionToken := generateSessionToken()
	expiration := time.Now().Add(ExpirationTime)

	session := Session{
		Username: username,
		Token:    sessionToken,
		Expires:  expiration,
	}
	db.mu.Lock()
	db.sessions[sessionToken] = session
	db.mu.Unlock()
	return &session
}

func (db SessionDB) CheckSession(r *http.Request) (*Session, bool) {
	cookie, err := r.Cookie(SessionToken)
	if err != nil {
		return nil, false
	}

	db.mu.RLock()
	session, exists := db.sessions[cookie.Value]
	db.mu.RUnlock()
	if !exists || session.Expires.Before(time.Now()) {
		return nil, false
	}
	return &session, true
}

func generateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (db *SessionDB) DeleteSession(token string) {
	db.mu.Lock()
	delete(db.sessions, token)
	db.mu.Unlock()
}
