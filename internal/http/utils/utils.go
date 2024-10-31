package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

const uploadPath = "./static/images"

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

func GetSessionFromContext(ctx context.Context) (models.Session, bool) {
	session, ok := ctx.Value(sessionKey).(models.Session)
	if !ok || session.Token == "" {
		return session, false
	}
	return session, true
}

func SetSessionInContext(ctx context.Context, session models.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
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

func GenerateFilename(header *multipart.FileHeader) error {
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		return models.ErrInvalidImage
	}
	token := base64.URLEncoding.EncodeToString(bytes)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	filename := fmt.Sprintf("%s_%d", token, timestamp)
	extension := getFileExtension(header.Filename)

	switch extension {
	case "jpeg", "jpg", "gif", "png":
	default:
		return models.ErrInvalidImageFormat
	}
	header.Filename = fmt.Sprintf("%s.%s", filename, extension)
	return nil
}

func getFileExtension(fileName string) string {
	parts := strings.Split(fileName, ".")
	extension := parts[1]

	if extension == "" {
		return ""
	}

	extension = strings.ToLower(strings.TrimSpace(extension))
	return extension
}
