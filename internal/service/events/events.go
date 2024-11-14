//go:generate mockgen -source ./events.go -destination=./mocks/events.go -package=mocks

package eventService

import (
	"context"
	"fmt"
	"strings"

	"kudago/internal/models"
)

type EventService struct {
	EventDB EventDB
	ImageDB ImageDB
}

type EventDB interface {
	GetUpcomingEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error)
	GetPastEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetEventsByCategory(ctx context.Context, categoryID int, paginationParams models.PaginationParams) ([]models.Event, error)
	GetEventsByUser(ctx context.Context, userID int, paginationParams models.PaginationParams) ([]models.Event, error)
	GetEventByID(ctx context.Context, ID int) (models.Event, error)
	AddEvent(ctx context.Context, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, ID int) error
	UpdateEvent(ctx context.Context, event models.Event) (models.Event, error)
	SearchEvents(ctx context.Context, params models.SearchParams, paginationParams models.PaginationParams) ([]models.Event, error)
	AddEventToFavorites(ctx context.Context, newFavorite models.FavoriteEvent) error
	DeleteEventFromFavorites(ctx context.Context, favorite models.FavoriteEvent) error
}

type ImageDB interface {
	SaveImage(ctx context.Context, media models.MediaFile) (string, error)
	DeleteImage(ctx context.Context, path string) error
}

func NewService(eventDB EventDB, imageDB ImageDB) EventService {
	return EventService{EventDB: eventDB, ImageDB: imageDB}
}

func (s *EventService) AddEvent(ctx context.Context, event models.Event, media models.MediaFile) (models.Event, error) {
	if media.File != nil {
		path, err := s.ImageDB.SaveImage(ctx, media)
		if err != nil {
			return models.Event{}, err
		}
		event.ImageURL = path
	}

	event, err := s.EventDB.AddEvent(ctx, event)
	if err != nil {
		if event.ImageURL != "" {
			s.ImageDB.DeleteImage(ctx, event.ImageURL)
		}
		return models.Event{}, err
	}
	return event, nil
}

func (s *EventService) DeleteEvent(ctx context.Context, ID, AuthorID int) error {
	dbEvent, err := s.EventDB.GetEventByID(ctx, ID)
	if err != nil {
		return err
	}

	if dbEvent.AuthorID != AuthorID {
		return fmt.Errorf("%s: %w", models.LevelService, models.ErrAccessDenied)
	}

	if dbEvent.ImageURL != "" {
		s.ImageDB.DeleteImage(ctx, dbEvent.ImageURL)
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

func (s *EventService) UpdateEvent(ctx context.Context, event models.Event, media models.MediaFile) (models.Event, error) {
	dbEvent, err := s.EventDB.GetEventByID(ctx, event.ID)
	if err != nil {
		return models.Event{}, err
	}

	if dbEvent.AuthorID != event.AuthorID {
		return models.Event{}, fmt.Errorf("%s: %w", models.LevelService, models.ErrAccessDenied)
	}

	if media.File != nil {
		path, err := s.ImageDB.SaveImage(ctx, media)
		if err != nil {
			return models.Event{}, err
		}
		event.ImageURL = path
	}

	updatedEvent, err := s.EventDB.UpdateEvent(ctx, event)
	if err != nil {
		if media.File != nil {
			s.ImageDB.DeleteImage(ctx, event.ImageURL)
		}
		return models.Event{}, err
	}

	if dbEvent.ImageURL != "" && media.File != nil {
		s.ImageDB.DeleteImage(ctx, dbEvent.ImageURL)
	}
	return updatedEvent, nil
}

func (s *EventService) AddEventToFavorites(ctx context.Context, newFavorite models.FavoriteEvent) error {
	return s.EventDB.AddEventToFavorites(ctx, newFavorite)
}

func (s *EventService) DeleteEventFromFavorites(ctx context.Context, favorite models.FavoriteEvent) error {
	return s.EventDB.DeleteEventFromFavorites(ctx, favorite)
}
