package eventService

import (
	"context"
	"strings"

	"kudago/internal/models"
)

type EventService struct {
	EventDB EventDB
}

type EventDB interface {
	GetAllEvents(ctx context.Context, offset, limit int) ([]models.Event, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetEventsByTags(ctx context.Context, tags []string) ([]models.Event, error)
	GetEventsByCategory(ctx context.Context, categoryID int) ([]models.Event, error)
	GetEventByID(ctx context.Context, ID int) (models.Event, error)
	AddEvent(ctx context.Context, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, ID int) error
	UpdateEvent(ctx context.Context, event models.Event) error
}

func NewService(eventDB EventDB) EventService {
	return EventService{EventDB: eventDB}
}

func (s *EventService) GetAllEvents(ctx context.Context, page, limit int) ([]models.Event, error) {
	offset := (page - 1) * limit
	return s.EventDB.GetAllEvents(ctx, offset, limit)
}

func (s *EventService) GetEventsByTags(ctx context.Context, tags []string) ([]models.Event, error) {
	for i, tag := range tags {
		tags[i] = strings.ToLower(tag)
	}
	return s.EventDB.GetEventsByTags(ctx, tags)
}

func (s *EventService) GetEventsByCategory(ctx context.Context, categoryID int) ([]models.Event, error) {
	return s.EventDB.GetEventsByCategory(ctx, categoryID)
}

func (s *EventService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.EventDB.GetCategories(ctx)
}

func (s *EventService) AddEvent(ctx context.Context, event models.Event) (models.Event, error) {
	return s.EventDB.AddEvent(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, ID, AuthorID int) error {
	dbEvent, err := s.EventDB.GetEventByID(ctx, ID)
	if err != nil {
		return err
	}

	if dbEvent.AuthorID != AuthorID {
		return models.ErrAccessDenied
	}
	return s.EventDB.DeleteEvent(ctx, ID)
}

func (s *EventService) GetEventByID(ctx context.Context, ID int) (models.Event, error) {
	return s.EventDB.GetEventByID(ctx, ID)
}

func (s *EventService) UpdateEvent(ctx context.Context, event models.Event) error {
	dbEvent, err := s.EventDB.GetEventByID(ctx, event.ID)
	if err != nil {
		return err
	}

	if dbEvent.AuthorID != event.AuthorID {
		return models.ErrAccessDenied
	}
	return s.EventDB.UpdateEvent(ctx, event)
}
