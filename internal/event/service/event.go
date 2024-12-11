package service

//go:generate mockgen -source ./event.go -destination=./mocks/events.go -package=mocks

import (
	"context"
	"fmt"
	"strings"

	"kudago/internal/models"
)

type EventService struct {
	EventDB EventDB
}

type EventDB interface {
	GetUpcomingEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error)
	GetPastEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetEventsByCategory(ctx context.Context, categoryID int, paginationParams models.PaginationParams) ([]models.Event, error)
	GetEventsByUser(ctx context.Context, userID int, paginationParams models.PaginationParams) ([]models.Event, error)
	GetEventByID(ctx context.Context, ID int) (models.Event, error)
	CreateEvent(ctx context.Context, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, ID int) error
	UpdateEvent(ctx context.Context, event models.Event) (models.Event, error)
	SearchEvents(ctx context.Context, params models.SearchParams, paginationParams models.PaginationParams) ([]models.Event, error)
	AddEventToFavorites(ctx context.Context, newFavorite models.FavoriteEvent) error
	DeleteEventFromFavorites(ctx context.Context, favorite models.FavoriteEvent) error
}

func NewService(eventDB EventDB) EventService {
	return EventService{EventDB: eventDB}
}

func (s *EventService) AddEvent(ctx context.Context, event models.Event) (models.Event, error) {
	return s.EventDB.CreateEvent(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, ID, AuthorID int) error {
	dbEvent, err := s.EventDB.GetEventByID(ctx, ID)
	if err != nil {
		return err
	}

	if dbEvent.AuthorID != AuthorID {
		return fmt.Errorf("%s: %w", models.LevelService, models.ErrAccessDenied)
	}

	return s.EventDB.DeleteEvent(ctx, ID)
}

func (s *EventService) GetEventByID(ctx context.Context, ID int) (models.Event, error) {
	return s.EventDB.GetEventByID(ctx, ID)
}

func (s *EventService) SearchEvents(ctx context.Context, params models.SearchParams, paginationParams models.PaginationParams) ([]models.Event, error) {
	for i, tag := range params.Tags {
		params.Tags[i] = strings.ToLower(tag)
	}
	return s.EventDB.SearchEvents(ctx, params, paginationParams)
}

func (s *EventService) UpdateEvent(ctx context.Context, event models.Event) (models.Event, error) {
	dbEvent, err := s.EventDB.GetEventByID(ctx, event.ID)
	if err != nil {
		return models.Event{}, err
	}

	if dbEvent.AuthorID != event.AuthorID {
		return models.Event{}, fmt.Errorf("%s: %w", models.LevelService, models.ErrAccessDenied)
	}

	updatedEvent, err := s.EventDB.UpdateEvent(ctx, event)
	if err != nil {
		return models.Event{}, err
	}

	return updatedEvent, nil
}

func (s *EventService) AddEventToFavorites(ctx context.Context, newFavorite models.FavoriteEvent) error {
	return s.EventDB.AddEventToFavorites(ctx, newFavorite)
}

func (s *EventService) DeleteEventFromFavorites(ctx context.Context, favorite models.FavoriteEvent) error {
	return s.EventDB.DeleteEventFromFavorites(ctx, favorite)
}
