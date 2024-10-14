package models

import "time"

type Session struct {
	Username string
	Token    string
	Expires  time.Time
}
