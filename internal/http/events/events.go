package events

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"kudago/internal/http/pkg"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type EventRequest struct {
	Title       string   `json:"title" valid:"required,length(3|50)"`
	Description string   `json:"description" valid:"required"`
	Tag         []string `json:"tag"`
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
}

type EventResponse struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tag         []string `json:"tag"`
	AuthorID    int      `json:"author"`
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
}

type CreateEventResponse struct {
	Event EventResponse `json:"event"`
}

type GetEventsResponse struct {
	Events []EventResponse `json:"events"`
}

type EventHandler struct {
	Service EventService
}

type EventService interface {
	GetAllEvents(ctx context.Context) []models.Event
	GetEventsByTag(ctx context.Context, tag string) []models.Event
	GetEventByID(ctx context.Context, ID int) (models.Event, error)
	AddEvent(ctx context.Context, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, ID int, authorID int) error
	UpdateEvent(ctx context.Context, event models.Event) error
}

func NewEventHandler(s EventService) *EventHandler {
	return &EventHandler{
		Service: s,
	}
}

func (h EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events := h.Service.GetAllEvents(r.Context())
	resp := GetEventsResponse{}
	for _, event := range events {
		eventResp := eventToEventResponse(event)
		resp.Events = append(resp.Events, eventResp)
	}
	json.NewEncoder(w).Encode(resp)
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
	resp := GetEventsResponse{}
	for _, event := range filteredEvents {
		eventResp := eventToEventResponse(event)
		resp.Events = append(resp.Events, eventResp)
	}
	json.NewEncoder(w).Encode(resp)
}

func (h EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	event, err := h.Service.GetEventByID(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			w.WriteHeader(http.StatusNoContent)
		default:
			pkg.WriteResponse(w, http.StatusInternalServerError, errInternal)
		}
		return
	}
	resp := eventToEventResponse(event)
	json.NewEncoder(w).Encode(resp)
}

func (h EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	sessionInfo, ok := pkg.GetSessionFromContext(r.Context())
	if !ok || !sessionInfo.Authenticated {
		pkg.WriteResponse(w, http.StatusForbidden, errUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorID:=sessionInfo.Session.UserID
	err = h.Service.DeleteEvent(r.Context(), id, authorID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			pkg.WriteResponse(w, http.StatusNotFound, errEventNotFound)
		case errors.Is(err, models.ErrAccessDenied):
			pkg.WriteResponse(w, http.StatusForbidden, errAccessDenied)
		default:
			pkg.WriteResponse(w, http.StatusInternalServerError, errInternal)
		}
		return
	}
}

func (h EventHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	sessionInfo, ok := pkg.GetSessionFromContext(r.Context())
	if !ok || !sessionInfo.Authenticated {
		pkg.WriteResponse(w, http.StatusForbidden, errUnauthorized)
		return
	}

	var req EventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, errInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		pkg.ProcessValidationErrors(w, err)
		return
	}

	event := models.Event{
		Title:       req.Title,
		Description: req.Description,
		DateStart:   req.DateStart,
		DateEnd:     req.DateEnd,
		AuthorID:    sessionInfo.Session.UserID,
		Tag:         req.Tag,
	}

	event, err = h.Service.AddEvent(r.Context(), event)

	if err != nil {
		switch {
		///TODO пока оставлю так, когда будет более четкая бд и ошибки для обработки, поправлю
		default:
			pkg.WriteResponse(w, http.StatusInternalServerError, errInternal)
		}
		return
	}
	resp := eventToEventResponse(event)
	json.NewEncoder(w).Encode(resp)
}

func (h EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	sessionInfo, ok := pkg.GetSessionFromContext(r.Context())
	if !ok || !sessionInfo.Authenticated {
		pkg.WriteResponse(w, http.StatusForbidden, errUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req EventRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, errInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		pkg.ProcessValidationErrors(w, err)
		return
	}

	event := models.Event{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		DateStart:   req.DateStart,
		DateEnd:     req.DateEnd,
		AuthorID:    sessionInfo.Session.UserID,
		Tag:         req.Tag,
	}

	err = h.Service.UpdateEvent(r.Context(), event)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			pkg.WriteResponse(w, http.StatusNotFound, errEventNotFound)
		case errors.Is(err, models.ErrAccessDenied):
			pkg.WriteResponse(w, http.StatusForbidden, errAccessDenied)
		default:
			pkg.WriteResponse(w, http.StatusInternalServerError, errInternal)
		}
		return
	}
	resp := eventToEventResponse(event)
	json.NewEncoder(w).Encode(resp)
}

func eventToEventResponse(event models.Event) EventResponse {
	return EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		DateStart:   event.DateStart,
		DateEnd:     event.DateEnd,
		Tag:         event.Tag,
	}
}
