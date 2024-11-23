package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	pbImage "kudago/internal/image/api"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

const (
	defaultPage   = 0
	defaultLimit  = 30
	maxUploadSize = 10 * 1024 * 1024 // 10Mb
)

func WriteResponse(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func HandleImageUpload(r *http.Request) (*pbImage.UploadRequest, error) {
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		return nil, models.ErrInvalidImage
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		if err != http.ErrMissingFile {
			return nil, models.ErrInvalidImage
		}
		return nil, nil
	}
	defer file.Close()

	err = GenerateFilename(header)
	if err != nil {
		return nil, models.ErrInvalidImage
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, models.ErrInvalidImage
	}

	return &pbImage.UploadRequest{
		Filename: header.Filename,
		File:     fileData,
	}, nil
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

type ValidationErrResponse struct {
	Errors []models.AuthError `json:"errors"`
}

type sessionKeyType struct{}

var sessionKey sessionKeyType

type requestIDKeyType struct{}

var requestIDKey requestIDKeyType

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
