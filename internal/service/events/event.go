package eventService

import (
	"context"
	"mime/multipart"
	"strings"

	"kudago/internal/models"
)

type EventService struct {
	EventDB EventDB
	ImageDB ImageDB
}

type EventDB interface {
	GetUpcomingEvents(ctx context.Context, offset, limit int) ([]models.Event, error)
	GetPastEvents(ctx context.Context, offset, limit int) ([]models.Event, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetEventsByTags(ctx context.Context, tags []string) ([]models.Event, error)
	GetEventsByCategory(ctx context.Context, categoryID int) ([]models.Event, error)
	GetEventsByUser(ctx context.Context, userID int) ([]models.Event, error)
	GetEventByID(ctx context.Context, ID int) (models.Event, error)
	AddEvent(ctx context.Context, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, ID int) error
	UpdateEvent(ctx context.Context, event models.Event) (models.Event, error)
	SearchEvents(ctx context.Context, paras models.SearchParams, offset, limit int) ([]models.Event, error)
}

type ImageDB interface {
	SaveImage(ctx context.Context, header multipart.FileHeader, file multipart.File) (string, error)
	DeleteImage(ctx context.Context, path string) error
}

func NewService(eventDB EventDB, imageDB ImageDB) EventService {
	return EventService{EventDB: eventDB, ImageDB: imageDB}
}

func (s *EventService) GetUpcomingEvents(ctx context.Context, page, limit int) ([]models.Event, error) {
	offset := (page - 1) * limit
	return s.EventDB.GetUpcomingEvents(ctx, offset, limit)
}

func (s *EventService) GetPastEvents(ctx context.Context, page, limit int) ([]models.Event, error) {
	offset := (page - 1) * limit
	return s.EventDB.GetPastEvents(ctx, offset, limit)
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

func (s *EventService) GetEventsByUser(ctx context.Context, userID int) ([]models.Event, error) {
	return s.EventDB.GetEventsByUser(ctx, userID)
}

func (s *EventService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.EventDB.GetCategories(ctx)
}

func (s *EventService) AddEvent(ctx context.Context, event models.Event, header *multipart.FileHeader, file *multipart.File) (models.Event, error) {
	path := ""
	if header != nil && file != nil {
		path, err := s.ImageDB.SaveImage(ctx, *header, *file)
		if err != nil {
			return models.Event{}, err
		}

		event.ImageURL = path
	}

	event, err := s.EventDB.AddEvent(ctx, event)
	if err != nil {
		if path != "" {
			s.ImageDB.DeleteImage(ctx, path)
		}
		return models.Event{}, err
	}
	return event, err
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

func (s *EventService) SearchEvents(ctx context.Context, params models.SearchParams, page, limit int) ([]models.Event, error) {
	offset := (page - 1) * limit

	for i, tag := range params.Tags {
		params.Tags[i] = strings.ToLower(tag)
	}
	return s.EventDB.SearchEvents(ctx, params, limit, offset)
}

func (s *EventService) UpdateEvent(ctx context.Context, event models.Event, header *multipart.FileHeader, file *multipart.File) (models.Event, error) {
	dbEvent, err := s.EventDB.GetEventByID(ctx, event.ID)
	if err != nil {
		return models.Event{}, err
	}

	if dbEvent.AuthorID != event.AuthorID {
		return models.Event{}, models.ErrAccessDenied
	}
	path := ""
	if header != nil && file != nil {
		path, err := s.ImageDB.SaveImage(ctx, *header, *file)
		if err != nil {
			return models.Event{}, err
		}

		event.ImageURL = path
	}

	event, err = s.EventDB.UpdateEvent(ctx, event)
	if err != nil {
		if path != "" {
			s.ImageDB.DeleteImage(ctx, path)
		}
		return models.Event{}, err
	}

	if dbEvent.ImageURL != "" && file != nil {
		err = s.ImageDB.DeleteImage(ctx, dbEvent.ImageURL)
	}
	return event, nil
}
