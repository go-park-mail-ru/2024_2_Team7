package models

import "errors"

var (
	ErrEmailIsUsed    = errors.New("email is already used")
	ErrUsernameIsUsed = errors.New("user alresdy exists")
	ErrEventNotFound  = errors.New("event not found")
	ErrAccessDenied   = errors.New("user has no access to event")
)
