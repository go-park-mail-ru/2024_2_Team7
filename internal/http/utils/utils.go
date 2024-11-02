package utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

type csrfKeyType struct{}

var csrfKey sessionKeyType

type requestIDKeyType struct{}

var requestIDKey requestIDKeyType

func WriteResponse(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func GetSessionFromContext(ctx context.Context) (models.Session, bool) {
	session, ok := ctx.Value(sessionKey).(models.Session)
	if !ok {
		return session, false
	}
	return session, true
}

func SetSessionInContext(ctx context.Context, session models.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

func GetCSRFFromContext(ctx context.Context) (models.TokenData, bool) {
	csrfKey, ok := ctx.Value(csrfKey).(models.TokenData)
	if !ok {
		return csrfKey, false
	}
	return csrfKey, true
}

func SetCSRFInContext(ctx context.Context, token models.TokenData) context.Context {
	return context.WithValue(ctx, csrfKey, token)
}

func GetRequestIDFromContext(ctx context.Context) string {
	ID, _ := ctx.Value(requestIDKey).(string)
	return ID
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
	WriteResponse(w, http.StatusBadRequest, resp)
}

func GetQueryParamInt(r *http.Request, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	value, err := strconv.Atoi(valueStr)

	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}

func LogRequestData(ctx context.Context, logger *zap.SugaredLogger, msg string, statusCode int, method, url, remoteAddr string, duration time.Duration, data map[string]interface{}) {
	requestID := GetRequestIDFromContext(ctx)

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
