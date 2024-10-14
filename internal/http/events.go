package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type EventHandler struct {
	Service EventService
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
