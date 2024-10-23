package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

type ValidationErrResponse struct {
	Errors []models.AuthError `json:"errors"`
}

type sessionKeyType struct{}

var sessionKey sessionKeyType

func WriteResponse(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func GetSessionFromContext(ctx context.Context) (*models.Session, bool) {
	session, ok := ctx.Value(sessionKey).(*models.Session)
	if session == nil {
		return session, false
	}
	return session, ok
}

func SetSessionInContext(ctx context.Context, session *models.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

func ProcessValidationErrors(w http.ResponseWriter, err error) {
	resp := ValidationErrResponse{}
	validationErrors := err.(govalidator.Errors)

	for _, err := range validationErrors {
		var validationErr govalidator.Error

		if errors.As(err, &validationErr) {
			valErr := models.AuthError{
				Field:   validationErr.Name,
				Message: validationErr.Err.Error(),
			}
			resp.Errors = append(resp.Errors, valErr)
		}
	}
	WriteResponse(w, http.StatusUnauthorized, resp)
}
