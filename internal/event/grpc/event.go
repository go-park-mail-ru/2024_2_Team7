package http

import (
	"context"

	pb "kudago/internal/event/api"
	"kudago/internal/logger"
	"kudago/internal/models"
)

const (
	errInternal           = "internal error"
	errEventNotFound      = "event not found"
	errAlreadyInFavorites = "event is already in favorites"
	errBadData            = "bad data request"
)

type ServerAPI struct {
	pb.UnimplementedEventServiceServer
	service EventService
	getter  EventsGetter
	logger  *logger.Logger
}

type EventService interface {
	AddEvent(ctx context.Context, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, ID, authorID int) error
	UpdateEvent(ctx context.Context, event models.Event) (models.Event, error)
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
	GetSubscriptionEvents(ctx context.Context, userID int, paginationParams models.PaginationParams) ([]models.Event, error)
}

func NewServerAPI(service EventService, getter EventsGetter, logger *logger.Logger) *ServerAPI {
	return &ServerAPI{
		service: service,
		getter:  getter,
		logger:  logger,
	}
}

func eventPBToEvent(event *pb.Event) models.Event {
	return models.Event{
		ID:          int(event.ID),
		Title:       event.Title,
		Description: event.Description,
		EventStart:  event.EventStart,
		EventEnd:    event.EventEnd,
		Location:    event.Location,
		Capacity:    int(event.Capacity),
		CategoryID:  int(event.CategoryID),
		Tag:         event.Tag,
		AuthorID:    int(event.AuthorID),
		ImageURL:    event.Image,
		Latitude:    float64(event.Latitude),
		Longitude:   float64(event.Longitude),
	}
}

func eventToEventPB(event models.Event) *pb.Event {
	return &pb.Event{
		ID:          int32(event.ID),
		Title:       event.Title,
		Description: event.Description,
		EventStart:  event.EventStart,
		EventEnd:    event.EventEnd,
		Location:    event.Location,
		Capacity:    int32(event.Capacity),
		CategoryID:  int32(event.CategoryID),
		Tag:         event.Tag,
		AuthorID:    int32(event.AuthorID),
		Image:       event.ImageURL,
		Latitude:    float64(event.Latitude),
		Longitude:   float64(event.Longitude),
	}
}

func favoritePBToFavorite(favorite *pb.FavoriteEvent) models.FavoriteEvent {
	return models.FavoriteEvent{
		EventID: int(favorite.EventID),
		UserID:  int(favorite.UserID),
	}
}

func writeEventsResponse(events []models.Event, limit int) *pb.Events {
	pbEvents := make([]*pb.Event, 0, limit)
	for _, event := range events {
		pbEvents = append(pbEvents, eventToEventPB(event))
	}

	return &pb.Events{
		Events: pbEvents,
	}
}

func getPaginationParams(params *pb.PaginationParams) models.PaginationParams {
	return models.PaginationParams{
		Limit:  int(params.Limit),
		Offset: int(params.Offset),
	}
}
