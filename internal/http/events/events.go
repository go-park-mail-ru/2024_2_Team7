package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	Location    string   `json:"location"`
	Category    int      `json:"category_id"`
	Capacity    int      `json:"capacity"`
	Tag         []string `json:"tag"`
	EventStart  string   `json:"event_start" valid:"rfc3339"`
	EventEnd    string   `json:"event_end" valid:"rfc3339"`
}

type EventResponse struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Category    int      `json:"category_id"`
	Capacity    int      `json:"capacity"`
	Tag         []string `json:"tag"`
	AuthorID    int      `json:"author"`
	EventStart  string   `json:"event_start"`
	EventEnd    string   `json:"event_end"`
	ImageURL    string   `json:"image"`
}

type CreateEventResponse struct {
	Event EventResponse `json:"event"`
}

type GetEventsResponse struct {
	Events []EventResponse `json:"events"`
}

type EventHandler struct {
	service EventService
}

type EventService interface {
	GetAllEvents(ctx context.Context, page, limit int) ([]models.Event, error)
	GetEventsByTags(ctx context.Context, tags []string) ([]models.Event, error)
	GetEventsByCategory(ctx context.Context, categoryID int) ([]models.Event, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetEventByID(ctx context.Context, ID int) (models.Event, error)
	AddEvent(ctx context.Context, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, ID, authorID int) error
	UpdateEvent(ctx context.Context, event models.Event) error
}

func NewEventHandler(s EventService) *EventHandler {
	return &EventHandler{
		service: s,
	}
}

// @Summary Получить все события
// @Description Получить все существующие события
// @Tags events
// @Accept  json
// @Produce  json
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество событий на странице (по умолчанию 30)"
// @Success 200 {object} GetEventsResponse
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events [get]
func (h EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	page := utils.GetQueryParamInt(r, "page", 1)
	limit := utils.GetQueryParamInt(r, "limit", 30)

	events, err := h.service.GetAllEvents(r.Context(), page, limit)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}
	resp := writeEventsResponse(events, limit)

	utils.WriteResponse(w, http.StatusOK, resp)
}

// пока просто ручка потом когда сделаем полноценный поиск поменяем

// @Summary Получение событий по тегу
// @Description Возвращает события по тегу
// @Tags events
// @Produce  json
// @Success 200 {object} GetEventsResponse
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/tags/{tag} [get]
func (h EventHandler) GetEventsByTags(w http.ResponseWriter, r *http.Request) {
	tagsParam := r.URL.Query().Get("tags")
	tags := strings.Split(tagsParam, ",")

	filteredEvents, err := h.service.GetEventsByTags(r.Context(), tags)
	if err != nil {
		fmt.Println(err)

		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := GetEventsResponse{}
	for _, event := range filteredEvents {
		eventResp := eventToEventResponse(event)
		resp.Events = append(resp.Events, eventResp)
	}
	utils.WriteResponse(w, http.StatusOK, resp)
}

// @Summary Получение событий по категори
// @Description Возвращает события по ID категории
// @Tags events
// @Produce  json
// @Success 200 {object} GetEventsResponse
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/categories/{category} [get]
func (h EventHandler) GetEventsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	categoryID, err := strconv.Atoi(category)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidCategory)
		return
	}

	filteredEvents, err := h.service.GetEventsByCategory(r.Context(), categoryID)
	if err != nil {
		switch err {
		case models.ErrInvalidCategory:
			utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidCategory)

		///TODO пока оставлю так, когда будет более четкая бд и ошибки для обработки, поправлю
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
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

	event, err := h.service.GetEventByID(r.Context(), id)
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
	err = h.service.DeleteEvent(r.Context(), id, authorID)
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
		fmt.Println(err)
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
		Location:    req.Location,
		EventStart:  req.EventStart,
		EventEnd:    req.EventEnd,
		AuthorID:    session.UserID,
		CategoryID:  req.Category,
		Capacity:    req.Capacity,
		Tag:         req.Tag,
	}

	event, err = h.service.AddEvent(r.Context(), event)
	if err != nil {
		switch err {
		case models.ErrInvalidCategory:
			utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidCategory)

		///TODO пока оставлю так, когда будет более четкая бд и ошибки для обработки, поправлю
		default:
			fmt.Println(err)
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
		EventStart:  req.EventStart,
		EventEnd:    req.EventEnd,
		AuthorID:    session.UserID,
		Tag:         req.Tag,
		Location:    req.Location,
		CategoryID:  req.Category,
		Capacity:    req.Capacity,
	}

	err = h.service.UpdateEvent(r.Context(), event)
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

// @Summary Получить все категории
// @Description Получить список всех доступных категорий событий
// @Tags categories
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Category "Список категорий"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /categories [get]
func (h EventHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories(r.Context())
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	utils.WriteResponse(w, http.StatusOK, categories)
}

func eventToEventResponse(event models.Event) EventResponse {
	return EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		EventStart:  event.EventStart,
		EventEnd:    event.EventEnd,
		Tag:         event.Tag,
		AuthorID:    event.AuthorID,
		Category:    event.CategoryID,
		ImageURL:    event.ImageURL,
		Capacity:    event.Capacity,
	}
}

func writeEventsResponse(events []models.Event, limit int) GetEventsResponse {
	resp := GetEventsResponse{make([]EventResponse, 0, limit)}

	for _, event := range events {
		eventResp := eventToEventResponse(event)
		resp.Events = append(resp.Events, eventResp)
	}
	return resp
}
