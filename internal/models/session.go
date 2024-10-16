package models

import "time"

type Session struct {
	UserID  int
	Token   string
	Expires time.Time
}
