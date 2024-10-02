package events

import "sync"

type Event struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
	Tag         []string `json:"tag"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image"`
}

type EventDB struct {
	Events []Event
	mu     *sync.RWMutex
}

func (db EventDB) GetAllEvents() []Event {
	db.mu.RLock()
	events := db.Events
	db.mu.RUnlock()
	return events
}

func (db EventDB) GetEventsByTag(tag string) []Event {
	var filteredEvents []Event
	db.mu.RLock()
	for _, event := range db.Events {
		for _, eventTag := range event.Tag {
			if tag == eventTag {
				filteredEvents = append(filteredEvents, event)
			}
		}
	}
	db.mu.RUnlock()
	return filteredEvents
}
