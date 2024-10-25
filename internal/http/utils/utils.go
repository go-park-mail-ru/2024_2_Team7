package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

type ValidationErrResponse struct {
	Errors []models.AuthError `json:"errors"`
}

type sessionKeyType struct{}

var sessionKey sessionKeyType

type requestIDKeyType struct{}

var requestIDKey requestIDKeyType

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

func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	ID, ok := ctx.Value(requestIDKey).(string)
	return ID, ok
}

func SetRequestIDInContext(ctx context.Context, ID string) context.Context {
	return context.WithValue(ctx, requestIDKey, ID)
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

func LogRequestData(ctx context.Context, logger *zap.SugaredLogger, msg string, statusCode int, method, url, remoteAddr string, duration time.Duration, data map[string]interface{}) {
	requestID, ok := GetRequestIDFromContext(ctx)
	if !ok {
		requestID = "unknown"
	}

	if data != nil {
		logger.Infow(msg,
			"request_id", requestID,
			"method", method,
			"url", url,
			"remote_addr", remoteAddr,
			"status_code", statusCode,
			"work_time", duration,
			"data", data,
		)
	} else {
		logger.Infow(msg,
			"request_id", requestID,
			"method", method,
			"url", url,
			"remote_addr", remoteAddr,
			"status_code", statusCode,
			"work_time", duration,
		)
	}
}
