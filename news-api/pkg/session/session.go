package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

type Session struct {
	Username string
	Expires  time.Time
}

var sessions = make(map[string]Session)

func CreateSession(username string) string {
	sessionToken := generateSessionToken()
	expiration := time.Now().Add(30 * time.Minute)
	sessions[sessionToken] = Session{
		Username: username,
		Expires:  expiration,
	}
	return sessionToken
}

func CheckSession(r *http.Request) (*Session, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, false
	}
	session, exists := sessions[cookie.Value]
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
func DeleteSession(token string) {
	delete(sessions, token)
}

func GetExpiration() time.Time {
	return time.Now().Add(30 * time.Minute)
}
