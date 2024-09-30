package events

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func NewEventDB() *EventDB {
	eventsFeed := createEventMapWithDefaultValues()

	return &EventDB{
		Events: eventsFeed,
	}
}

func (h *Handler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events:=h.EventDB.GetAllEvents()
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) GetEventsByTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["tag"]
	tag = strings.ToLower(tag)

	filteredEvents:=h.EventDB.GetEventsByTag(tag)

	if len(filteredEvents) == 0 {
        w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filteredEvents)
}

func createEventMapWithDefaultValues() []Event {
	events := []Event{
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

func (db* EventDB) GetAllEvents() []Event{
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.Events
}

func (db* EventDB) GetEventsByTag(tag string) []Event{
	db.mu.RLock()
	defer db.mu.RUnlock()

	var filteredEvents []Event
	for _, event := range db.Events {
		for _, eventTag := range event.Tag {
			if tag == eventTag {
				filteredEvents = append(filteredEvents, event)
			}
		}
	}
	return filteredEvents
}