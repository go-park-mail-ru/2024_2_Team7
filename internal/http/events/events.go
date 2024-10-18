package events

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"

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

// @Summary Получить все события
// @Description Получить все существующие события
// @Tags events
// @Accept  json
// @Produce  json
// @Success 200 {object} GetEventsResponse
// @Router /events [get]
func (h EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events := h.Service.GetAllEvents(r.Context())
	resp := GetEventsResponse{}

	for _, event := range events {
		eventResp := eventToEventResponse(event)
		resp.Events = append(resp.Events, eventResp)
	}
	utils.WriteResponse(w, http.StatusOK, resp)
}

// @Summary Получение событий по тегу
// @Description Возвращает события по тегу
// @Tags events
// @Produce  json
// @Success 200 {object} GetEventsResponse
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/{tag} [get]
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
	utils.WriteResponse(w, http.StatusOK, resp)
}

// @Summary Получение события по ID
// @Description Возвращает информацию о событии по его идентификатору
// @Tags events
// @Produce  json
// @Success 200 {object} EventResponse
// @Failure 404 {object} httpErrors.HttpError "Event Not Found"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/{id} [get]
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
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}
	resp := eventToEventResponse(event)
	utils.WriteResponse(w, http.StatusOK, resp)
}

// @Summary Удаление события
// @Description Удаляет существующее событие
// @Tags events
// @Produce  json
// @Success 204
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Failure 403 {object} httpErrors.HttpError "Access Denied"
// @Failure 404 {object} httpErrors.HttpError "Event Not Found"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/{id} [delete]
func (h EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())

	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorID := session.UserID
	err = h.Service.DeleteEvent(r.Context(), id, authorID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrEventNotFound)
		case errors.Is(err, models.ErrAccessDenied):
			utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrAccessDenied)
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}
}

// @Summary Создание события
// @Description Создает новое событие в системе
// @Tags events
// @Accept  json
// @Produce  json
// @Success 201 {object} EventResponse
// @Failure 400 {object} httpErrors.HttpError "Invalid Data"
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events [post]
func (h EventHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	var req EventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	event := models.Event{
		Title:       req.Title,
		Description: req.Description,
		DateStart:   req.DateStart,
		DateEnd:     req.DateEnd,
		AuthorID:    session.UserID,
		Tag:         req.Tag,
	}

	event, err = h.Service.AddEvent(r.Context(), event)
	if err != nil {
		switch {
		///TODO пока оставлю так, когда будет более четкая бд и ошибки для обработки, поправлю
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}
	resp := eventToEventResponse(event)
	utils.WriteResponse(w, http.StatusOK, resp)
}

// @Summary Обновление события
// @Description Обновляет данные существующего события
// @Tags events
// @Accept  json
// @Produce  json
// @Success 200 {object} EventResponse
// @Failure 400 {object} httpErrors.HttpError "Invalid Data"
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Failure 403 {object} httpErrors.HttpError "Access Denied"
// @Failure 404 {object} httpErrors.HttpError "Event Not Found"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/{id} [put]
func (h EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
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
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		utils.ProcessValidationErrors(w, err)
		return
	}

	event := models.Event{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		DateStart:   req.DateStart,
		DateEnd:     req.DateEnd,
		AuthorID:    session.UserID,
		Tag:         req.Tag,
	}

	err = h.Service.UpdateEvent(r.Context(), event)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrEventNotFound)
		case errors.Is(err, models.ErrAccessDenied):
			utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrAccessDenied)
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}
	resp := eventToEventResponse(event)
	utils.WriteResponse(w, http.StatusOK, resp)
}

func eventToEventResponse(event models.Event) EventResponse {
	return EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		DateStart:   event.DateStart,
		DateEnd:     event.DateEnd,
		Tag:         event.Tag,
		AuthorID:    event.AuthorID,
	}
}
