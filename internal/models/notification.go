//go:generate easyjson notification.go
package models

import "time"

//easyjson:json
type Notification struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	EventID  int       `json:"event_id"`
	NotifyAt time.Time `json:"notify_at"`
	Message  string    `json:"message"`
}
