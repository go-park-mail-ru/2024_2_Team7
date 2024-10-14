package http

import (
	"context"

	"kudago/internal/models"
)

type AuthService interface {
	CheckSession(ctx context.Context, cookie string) (*models.Session, bool)
	UserExists(ctx context.Context, username string) bool
	GetUser(ctx context.Context, username string) models.User
	CheckCredentials(ctx context.Context, creds models.Credentials) bool
	Register(ctx context.Context, user models.User) error
	CreateSession(ctx context.Context, username string) *models.Session
	DeleteSession(ctx context.Context, username string)
}

type EventService interface {
	GetAllEvents(ctx context.Context) []models.Event
	GetEventsByTag(ctx context.Context, tag string) []models.Event
}
