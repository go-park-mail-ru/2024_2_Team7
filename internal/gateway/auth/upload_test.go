package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/gateway/user/mocks"
	pbImage "kudago/internal/image/api"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandlers_UploadImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageService := mocks.NewMockImageServiceClient(ctrl)
	logger, _ := logger.NewLogger()

	handler := &AuthHandlers{
		ImageService: mockImageService,
		logger:       logger,
	}

	tests := []struct {
		name      string
		req       *pbImage.UploadRequest
		wantCode  int
		wantError error
	}{
		{
			name: "Успешная загрузка изображения",
			req: &pbImage.UploadRequest{
				Filename: "test_image.png",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Неверный формат изображения",
			req: &pbImage.UploadRequest{
				Filename: "invalid_image.txt",
			},
			wantCode:  http.StatusBadRequest,
			wantError: models.ErrInvalidImageFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			if tt.wantError == nil {
				mockImageService.EXPECT().UploadImage(gomock.Any(), tt.req).Return(&pbImage.UploadResponse{FileUrl: "http://example.com/image.png"}, nil)
			} else {
				mockImageService.EXPECT().UploadImage(gomock.Any(), tt.req).Return(nil, tt.wantError)
			}

			url, err := handler.uploadImage(context.Background(), tt.req, recorder)

			assert.Equal(t, tt.wantCode, recorder.Code)
			if tt.wantError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "http://example.com/image.png", url)
			}
		})
	}
}

func TestAuthHandlers_DeleteImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageService := mocks.NewMockImageServiceClient(ctrl)
	logger, _ := logger.NewLogger()

	handler := &AuthHandlers{
		ImageService: mockImageService,
		logger:       logger,
	}

	tests := []struct {
		name     string
		url      string
		wantCode int
	}{
		{
			name:     "Пустой URL",
			url:      "",
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			if tt.url != "" {
				mockImageService.EXPECT().DeleteImage(gomock.Any(), &pbImage.DeleteRequest{FileUrl: tt.url}).Return(nil) // Добавьте обработку ошибок при необходимости
			}

			handler.deleteImage(context.Background(), tt.url)

			assert.Equal(t, tt.wantCode, recorder.Code)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})
	}
}
