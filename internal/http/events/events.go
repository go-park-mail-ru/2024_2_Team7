//go:generate mockgen -source ./events.go -destination=./mocks/events.go -package=mocks

package events

import (
	"context"

	"kudago/internal/logger"

	"kudago/internal/models"
)

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
}

type EventsGetter interface {
	GetUpcomingEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error)
	GetPastEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error)
	GetEventsByCategory(ctx context.Context, categoryID int, paginationParams models.PaginationParams) ([]models.Event, error)
	GetEventsByUser(ctx context.Context, userID int, paginationParams models.PaginationParams) ([]models.Event, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetEventByID(ctx context.Context, ID int) (models.Event, error)
}

func NewEventHandler(s EventService, g EventsGetter, logger *logger.Logger) *EventHandler {
	return &EventHandler{
		service: s,
		logger:  logger,
		getter:  g,
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
