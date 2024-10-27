package models

import "time"

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventStart  string    `json:"event_start"`
	EventEnd    string    `json:"event_finish"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
	CategoryID  int       `json:"category_id"`
	AuthorID    int       `json:"author"`
	Tag         []string  `json:"tag"`
	ImageURL    string    `json:"image"`
}
