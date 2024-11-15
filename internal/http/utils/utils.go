package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

const (
	uploadPath    = "./static/images"
	defaultPage   = 0
	defaultLimit  = 30
	maxUploadSize = 1 * 1024 * 1024 // 1Mb
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
	extension := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileName), "."))
	return extension
}

func GetPaginationParams(r *http.Request) models.PaginationParams {
	page := GetQueryParamInt(r, "page", defaultPage)
	limit := GetQueryParamInt(r, "limit", defaultLimit)
	offset := page * limit
	return models.PaginationParams{
		Offset: offset,
		Limit:  limit,
	}
}

func HandleImageUpload(r *http.Request) (models.MediaFile, error) {
	r.ParseMultipartForm(maxUploadSize)
	file, header, err := r.FormFile("image")
	if err != nil {
		if err != http.ErrMissingFile {
			return models.MediaFile{}, models.ErrInvalidImage
		}
		return models.MediaFile{}, nil
	}

	err = GenerateFilename(header)
	if err != nil {
		return models.MediaFile{}, models.ErrInvalidImage
	}

	return models.MediaFile{
		Filename: header.Filename,
		File:     file,
	}, nil
}

// SanitizeStruct чистит всех строковых полей структуры от XSS
func SanitizeStruct(input interface{}) error {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input is not a struct")
	}

	p := bluemonday.UGCPolicy()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			original := field.String()
			field.SetString(p.Sanitize(original))
		}
	}
	return nil
}
