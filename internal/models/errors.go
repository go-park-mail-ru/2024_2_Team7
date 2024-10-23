package models

import "errors"

var (
	ErrEventNotFound = errors.New("event not found")
	ErrAccessDenied  = errors.New("user has no access to event")
	ErrUserNotFound  = errors.New("user not found")
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
