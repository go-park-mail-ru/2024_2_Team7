package eventRepository

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"kudago/internal/models"
)

type EventDB struct {
	Events []models.Event
	mu     *sync.RWMutex
}

func NewDB() *EventDB {
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

func (db EventDB) GetEventByID(ctx context.Context, ID int) (models.Event, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, event := range db.Events {
		if event.ID == ID {
			return event, nil
		}
	}

	return models.Event{}, models.ErrEventNotFound
}

func (db *EventDB) DeleteEvent(ctx context.Context, ID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	for i, event := range db.Events {
		if event.ID == ID {
			db.Events = slices.Delete(db.Events, i, i+1)
			return nil
		}
	}

	return models.ErrEventNotFound
}

func (db *EventDB) UpdateEvent(ctx context.Context, updatedEvent models.Event) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	for i, event := range db.Events {
		if event.ID == updatedEvent.ID {
			fmt.Println(event, updatedEvent)
			db.Events[i] = updatedEvent
			return nil
		}
	}

	return models.ErrEventNotFound
}

func (db *EventDB) AddEvent(ctx context.Context, event models.Event) (models.Event, error) {
	db.mu.Lock()
	event.ID = len(db.Events)
	db.Events = append(db.Events, event)
	db.mu.Unlock()
	return event, nil
}

func createEventMapWithDefaultValues() []models.Event {
	events := []models.Event{
		{
			ID:          0,
			Title:       "Выставка в Третьяковской галерее",
			DateStart:   "2024-09-27",
			DateEnd:     "2024-10-3",
			AuthorID:    1,
			Tag:         []string{"popular"},
			Description: "Творчество Малевича",
			ImageURL:    "/static/images/Event1.jpeg",
		},
		{
			ID:          1,
			Title:       "Экскурсии по центру",
			DateStart:   "2024-09-26",
			DateEnd:     "2024-11-5",
			AuthorID:    1,
			Tag:         []string{"popular"},
			Description: "Прогулка по Кремлю",
			ImageURL:    "/static/images/Event2.jpg",
		},
		{
			ID:          2,
			Title:       "Концерт в Концертном зале",
			DateStart:   "2024-10-15",
			DateEnd:     "2024-10-15",
			AuthorID:    2,
			Tag:         []string{"popular"},
			Description: "Музыкальное выступление известной группы",
			ImageURL:    "/static/images/PetrovConcert.jpg",
		},
		{
			ID:          3,
			Title:       "Фестиваль искусств",
			DateStart:   "2024-12-01",
			DateEnd:     "2024-12-05",
			AuthorID:    0,
			Tag:         []string{"festival"},
			Description: "Мультимедиальный фестиваль современного искусства",
			ImageURL:    "/static/images/SemenovArtFestival.jpg",
		},
	}
	return events
}
