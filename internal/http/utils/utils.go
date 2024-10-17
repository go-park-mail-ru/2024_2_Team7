package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

type ValidationErrResponse struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func WriteResponse(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func GetSessionFromContext(ctx context.Context) (*models.Session, bool) {
	session, ok := ctx.Value(models.SessionKey).(*models.Session)
	return session, ok
}

func ProcessValidationErrors(w http.ResponseWriter, err error) {
	resp := ValidationErrResponse{}
	errors := err.(govalidator.Errors)

	for _, err := range errors {
		if validationErr, ok := err.(govalidator.Error); ok {
			valErr := ValidationError{
				Field: validationErr.Name,
				Error: validationErr.Validator,
			}
			resp.Errors = append(resp.Errors, valErr)
		}
	}
	WriteResponse(w, http.StatusUnauthorized, resp)
}
