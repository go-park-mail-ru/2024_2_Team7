package models

import "errors"

var (
	ErrEmailIsUsed       = errors.New("Email is already used")
	ErrUserAlreadyExists = errors.New("User alresdy exists")
)
