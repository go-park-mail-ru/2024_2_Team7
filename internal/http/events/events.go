//go:generate mockgen -source ./events.go -destination=./mocks/events.go -package=mocks

package events

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/logger"

	"kudago/internal/models"
)

var maxDate = time.Date(2030, 12, 31, 0, 0, 0, 0, time.UTC)

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

type GetEventsResponse struct {
	Events []EventResponse `json:"events"`
}

type NewEventRequest struct {
	Title       string   `json:"title" valid:"required,length(3|100)"`
	Description string   `json:"description" valid:"required,length(3|1000)" `
	Location    string   `json:"location" valid:"length(3|100)"`
	Category    int      `json:"category_id" valid:"required,range(1|8)"`
	Capacity    int      `json:"capacity" valid:"range(1|20000)"`
	Tag         []string `json:"tag"`
	EventStart  string   `json:"event_start" valid:"rfc3339,required"`
	EventEnd    string   `json:"event_end" valid:"rfc3339,required"`
}

type NewEventResponse struct {
	Event EventResponse `json:"event"`
}

type EventHandler struct {
	service EventService
	getter  EventsGetter
	logger  *logger.Logger
}

type EventService interface {
	AddEvent(ctx context.Context, event models.Event, media models.MediaFile) (models.Event, error)
	DeleteEvent(ctx context.Context, ID, authorID int) error
	UpdateEvent(ctx context.Context, event models.Event, media models.MediaFile) (models.Event, error)
	SearchEvents(ctx context.Context, params models.SearchParams, paginationParams models.PaginationParams) ([]models.Event, error)
	AddEventToFavorites(ctx context.Context, newFavorite models.FavoriteEvent) error
	DeleteEventFromFavorites(ctx context.Context, newFavorite models.FavoriteEvent) error
}

type EventsGetter interface {
	GetUpcomingEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error)
	GetPastEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error)
	GetEventsByCategory(ctx context.Context, categoryID int, paginationParams models.PaginationParams) ([]models.Event, error)
	GetEventsByUser(ctx context.Context, userID int, paginationParams models.PaginationParams) ([]models.Event, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetEventByID(ctx context.Context, ID int) (models.Event, error)
	GetFavorites(ctx context.Context, userID int, paginationParams models.PaginationParams) ([]models.Event, error)
}

func NewEventHandler(s EventService, g EventsGetter, logger *logger.Logger) *EventHandler {
	return &EventHandler{
		service: s,
		logger:  logger,
		getter:  g,
	}
}

func checkNewEventRequest(req NewEventRequest) *httpErrors.HttpError {
	if len(req.Tag) > 3 {
		return httpErrors.ErrTooManyTags
	}

	for _, tag := range req.Tag {
		if len(tag) > 20 || len(tag) == 0 {
			return httpErrors.ErrBadTagLength
		}
	}

	eventStart, err := time.Parse(time.RFC3339, req.EventStart)
	if err != nil {
		return httpErrors.ErrInvalidTime
	}

	eventEnd, err := time.Parse(time.RFC3339, req.EventEnd)
	if err != nil {
		return httpErrors.ErrInvalidTime
	}

	if !eventEnd.After(eventStart) {
		return httpErrors.ErrEventStartAfterEventEnd
	}

	if eventStart.Before(time.Now()) || eventEnd.After(maxDate) {
		return httpErrors.ErrBadEventTiming
	}

	return nil
}

func parseEventData(r *http.Request) (NewEventRequest, models.MediaFile, *httpErrors.HttpError) {
	var req NewEventRequest
	var media models.MediaFile
	jsonData := r.FormValue("json")
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		return req, media, httpErrors.ErrInvalidData
	}

	media, err = utils.HandleImageUpload(r)
	if err != nil {
		return req, media, httpErrors.ErrInvalidImage
	}

	return req, media, nil
}

func toModelEvent(req NewEventRequest, authorID int) models.Event {
	return models.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		EventStart:  req.EventStart,
		EventEnd:    req.EventEnd,
		AuthorID:    authorID,
		CategoryID:  req.Category,
		Capacity:    req.Capacity,
		Tag:         req.Tag,
	}
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
