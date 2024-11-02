package models

import "errors"

var (
	ErrEventNotFound      = errors.New("event not found")
	ErrAccessDenied       = errors.New("user has no access to event")
	ErrUserNotFound       = errors.New("user not found")
	ErrInternal           = errors.New("internal error")
	ErrInvalidCategory    = errors.New("invalid category")
	ErrInvalidImageFormat = errors.New("invalid image format")
	ErrInvalidImage       = errors.New("invalid image")
	ErrUnsupportedFile    = errors.New("unsupported file type")
)

const (
	LevelDB      string = "DB"
	LevelService string = "Service"
	LevelHandler string = "Handler"
)

type AuthError struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

var (
	ErrEmailIsUsed = &AuthError{
		Field:   "email",
		Message: "email is already used",
	}

	ErrUsernameIsUsed = &AuthError{
		Field:   "username",
		Message: "user already exists",
	}
)

func (e AuthError) Error() string {
	return e.Message
}
