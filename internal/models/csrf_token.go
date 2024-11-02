package models

import "time"

type TokenData struct {
	CSRFtoken    string
	SessionToken string
	UserID       int
	Exp          time.Time
}
