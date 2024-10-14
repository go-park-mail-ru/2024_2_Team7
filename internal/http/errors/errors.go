package httpErrors

import (
	"net/http"
)

type AuthError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	ErrUserIsAuthorized = &AuthError{
		Message: "User is authorized",
		Code:    http.StatusForbidden,
	}

	ErrInvalidRequest= &AuthError{
		Message: "Invalid request",
		Code:    http.StatusBadRequest,
	}

	ErrUserAlreadyExists= &AuthError{
		Message: "User alresdy exists",
		Code:    http.StatusConflict,
	}

	ErrUserAlreadyLoggedIn= &AuthError{
		Message: "Already logged in",
		Code:    http.StatusForbidden,
	}

	ErrUnauthorized= &AuthError{
		Message: "Unauthorized",
		Code:    http.StatusUnauthorized,
	}
)
