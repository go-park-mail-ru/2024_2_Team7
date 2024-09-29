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

func (h *EventHandler) GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.EventDB.Events)
}

func (h *EventHandler) GetEventsByTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["tag"]
	tag = strings.ToLower(tag)

	var filteredEvents []Event
	for _, event := range h.EventDB.Events {
		for _, eventTag := range event.Tag {
			if tag == eventTag {
				filteredEvents = append(filteredEvents, event)
			}
		}
	}

	if len(filteredEvents) == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(filteredEvents)
}

func createEventMapWithDefaultValues() []Event {
	events := []Event{
		{
			ID:          "1",
			Title:       "Выставка в Третьяковской галерее",
			DateStart:   "2024-09-27",
			DateEnd:     "2024-10-3",
			Tag:         []string{"popular"},
			Description: "Творчество Малевича",
			ImageURL:    "/static/images/Events1.jpg",
		},
		{
			ID:          "2",
			Title:       "Экскурсии по центру",
			DateStart:   "2024-09-26",
			DateEnd:     "2024-11-5",
			Tag:         []string{"popular"},
			Description: "Прогулка по Кремлю",
			ImageURL:    "/static/images/Events2.jpg",
		},
		{
			ID:          "3",
			Title:       "Концерт в Концертном зале",
			DateStart:   "2024-10-15",
			DateEnd:     "2024-10-15",
			Tag:         []string{"popular"},
			Description: "Музыкальное выступление известной группы",
			ImageURL:    "/static/images/PetrovConcert.jpg",
		},
		{
			ID:          "4",
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
