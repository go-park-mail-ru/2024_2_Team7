package events

import "sync"

type Event struct {
	ID          int   `json:"id"`
	Title       string   `json:"title"`
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
	Tag         []string `json:"tag"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image"`
}

type EventDB struct {
	Events []Event
	mu     sync.RWMutex
}

type Handler struct {
	EventDB EventDB
}
