package models

import "time"
const(
	SessionToken = "session_token"
	SessionKey = "session"
)

type Session struct {
	UserID  int
	Token   string
	Expires time.Time
}

type SessionInfo struct {
	Session       Session
	Authenticated bool
}
