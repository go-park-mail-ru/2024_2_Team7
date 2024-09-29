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
	ExpirationTime = 30 * time.Minute
)

type Session struct {
	Username string
	Token    string
	Expires  time.Time
}

type SessionDB struct {
	sessions map[string]Session
	mu       sync.RWMutex
}

func NewSessionDB() *SessionDB {
	return &SessionDB{
		sessions: make(map[string]Session),
	}
}

func (db *SessionDB) CreateSession(username string) *Session {
	db.mu.Lock()
	defer db.mu.Unlock()

	sessionToken := generateSessionToken()
	expiration := time.Now().Add(ExpirationTime)

	session := Session{
		Username: username,
		Token:    sessionToken,
		Expires:  expiration,
	}

	db.sessions[sessionToken] = session
	return &session
}

func (db *SessionDB) CheckSession(r *http.Request) (*Session, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	cookie, err := r.Cookie(SessionToken)
	if err != nil {
		return nil, false
	}

	session, exists := db.sessions[cookie.Value]
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
	defer db.mu.Unlock()
	delete(db.sessions, token)
}
