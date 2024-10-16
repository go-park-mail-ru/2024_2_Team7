package models

import "errors"

var (
	ErrEmailIsUsed    = errors.New("Email is already used")
	ErrUsernameIsUsed = errors.New("User alresdy exists")
)
