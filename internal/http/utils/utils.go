package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"kudago/internal/models"
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

func GetSessionFromContext(ctx context.Context) (models.SessionInfo, bool) {
	sessionInfo, ok := ctx.Value(models.SessionKey).(models.SessionInfo)
	return sessionInfo, ok
}

func ProcessValidationErrors(w http.ResponseWriter, err error) {
	errors := strings.Split(err.Error(), ";")
	resp := ValidationErrResponse{}

	for _, err := range errors {
		colonIndex := strings.Index(err, ":")

		if colonIndex == -1 {
			continue
		}

		field := err[:colonIndex]
		errorMsg := err[colonIndex+2:]

		valErr := ValidationError{
			Field: field,
			Error: errorMsg,
		}
		resp.Errors = append(resp.Errors, valErr)
	}
	WriteResponse(w, http.StatusUnauthorized, resp)
}
