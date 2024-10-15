package events

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"kudago/internal/models"
)

type EventHandler struct {
	Service EventService
}

type EventService interface {
	GetAllEvents(ctx context.Context) []models.Event
	GetEventsByTag(ctx context.Context, tag string) []models.Event
}

func NewEventHandler(s EventService) *EventHandler {
	return &EventHandler{
		Service: s,
	}
}

func (h EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events := h.Service.GetAllEvents(r.Context())
	json.NewEncoder(w).Encode(events)
}

func (h EventHandler) GetEventsByTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["tag"]
	tag = strings.ToLower(tag)

	filteredEvents := h.Service.GetEventsByTag(r.Context(), tag)

	if len(filteredEvents) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filteredEvents)
}
