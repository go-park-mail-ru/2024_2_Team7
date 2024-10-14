package repository

import (
	"context"
	"sync"

	"kudago/internal/models"
)

type EventDB struct {
	Events []models.Event
	mu     *sync.RWMutex
}

func NewEventDB() *EventDB {
	eventsFeed := createEventMapWithDefaultValues()

	return &EventDB{
		Events: eventsFeed,
		mu:     &sync.RWMutex{},
	}
}

func (db EventDB) GetAllEvents(ctx context.Context) []models.Event {
	db.mu.RLock()
	events := db.Events
	db.mu.RUnlock()
	return events
}

func (db EventDB) GetEventsByTag(ctx context.Context, tag string) []models.Event {
	var filteredEvents []models.Event
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

func createEventMapWithDefaultValues() []models.Event {
	events := []models.Event{
		{
			ID:          1,
			Title:       "Выставка в Третьяковской галерее",
			DateStart:   "2024-09-27",
			DateEnd:     "2024-10-3",
			Tag:         []string{"popular"},
			Description: "Творчество Малевича",
			ImageURL:    "/static/images/Event1.jpeg",
		},
		{
			ID:          2,
			Title:       "Экскурсии по центру",
			DateStart:   "2024-09-26",
			DateEnd:     "2024-11-5",
			Tag:         []string{"popular"},
			Description: "Прогулка по Кремлю",
			ImageURL:    "/static/images/Event2.jpg",
		},
		{
			ID:          3,
			Title:       "Концерт в Концертном зале",
			DateStart:   "2024-10-15",
			DateEnd:     "2024-10-15",
			Tag:         []string{"popular"},
			Description: "Музыкальное выступление известной группы",
			ImageURL:    "/static/images/PetrovConcert.jpg",
		},
		{
			ID:          4,
			Title:       "Фестиваль искусств",
			DateStart:   "2024-12-01",
			DateEnd:     "2024-12-05",
			Tag:         []string{"festival"},
			Description: "Мультимедиальный фестиваль современного искусства",
			ImageURL:    "/static/images/SemenovArtFestival.jpg",
		},
	}
	return events
}
