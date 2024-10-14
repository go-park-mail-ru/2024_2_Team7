package service

import (
	"context"

	"kudago/internal/models"
	"kudago/internal/repository"
)

type EventService struct {
	EventDB iEventDB
}

func NewEventService(eventDB *repository.EventDB) EventService {
	return EventService{EventDB: eventDB}
}

func (s *EventService) GetAllEvents(ctx context.Context) []models.Event {
	return s.EventDB.GetAllEvents(ctx)
}

func (s *EventService) GetEventsByTag(ctx context.Context, tag string) []models.Event {
	return s.EventDB.GetEventsByTag(ctx, tag)
}
