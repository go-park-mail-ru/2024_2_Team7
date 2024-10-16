package models

import "errors"

var (
	ErrEmailIsUsed    = errors.New("Email is already used")
	ErrUsernameIsUsed = errors.New("User alresdy exists")
	ErrEventNotFound  = errors.New("Event not found")
	ErrAccessDenied   = errors.New("User has no access to event")
)
