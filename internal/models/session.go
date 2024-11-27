package models

import "time"

const (
	SessionToken = "session_token"
)

type Session struct {
	UserID  int
	Token   string
	Expires time.Time
}
