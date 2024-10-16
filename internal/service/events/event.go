package eventService

import (
	"context"

	"kudago/internal/models"
)

type EventService struct {
	EventDB EventDB
}

type EventDB interface {
	GetAllEvents(ctx context.Context) []models.Event
	GetEventsByTag(ctx context.Context, tag string) []models.Event
}

func NewService(eventDB EventDB) EventService {
	return EventService{EventDB: eventDB}
}

func (s *EventService) GetAllEvents(ctx context.Context) []models.Event {
	return s.EventDB.GetAllEvents(ctx)
}

func (s *EventService) GetEventsByTag(ctx context.Context, tag string) []models.Event {
	return s.EventDB.GetEventsByTag(ctx, tag)
}
