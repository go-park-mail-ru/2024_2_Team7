package service

import (
	"context"

	"kudago/internal/models"
)

type iUserDB interface {
	UserExists(ctx context.Context, username string) bool
	AddUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, username string) models.User
	CheckCredentials(ctx context.Context, username string, password string) bool
}

type iSessionDB interface {
	CheckSession(ctx context.Context, cookie string) (*models.Session, bool)
	CreateSession(ctx context.Context, username string) *models.Session
	DeleteSession(ctx context.Context, username string)
}

type iEventDB interface {
	GetAllEvents(ctx context.Context) []models.Event
	GetEventsByTag(ctx context.Context, tag string) []models.Event
}
