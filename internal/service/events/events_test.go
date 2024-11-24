package eventService

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"kudago/internal/models"
	"kudago/internal/service/events/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEventService_AddEvent(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventDB := mocks.NewMockEventDB(ctrl)
	mockImageDB := mocks.NewMockImageDB(ctrl)
	service := NewService(mockEventDB, mockImageDB)

	mockFile := io.NopCloser(bytes.NewReader([]byte("updated image data")))

	testCases := []struct {
		name        string
		event       models.Event
		media       models.MediaFile
		setupMocks  func()
		expectedURL string
		expectError bool
	}{
		{
			name: "успешное сохранение",
			event: models.Event{
				Title: "Event 1",
			},
			media: models.MediaFile{
				Filename: "image.jpg",
				File:     mockFile,
			},
			setupMocks: func() {
				mockImageDB.EXPECT().SaveImage(gomock.Any(), gomock.Any()).Return("path/to/image.jpg", nil)
				mockEventDB.EXPECT().AddEvent(gomock.Any(), gomock.Any()).Return(models.Event{
					ID:       1,
					Title:    "Event 1",
					ImageURL: "path/to/image.jpg",
				}, nil)
			},
			expectedURL: "path/to/image.jpg",
			expectError: false,
		},
		{
			name: "ошибка при сохранении",
			event: models.Event{
				Title: "Event 2",
			},
			media: models.MediaFile{
				Filename: "image.jpg",
				File:     mockFile,
			},
			setupMocks: func() {
				mockImageDB.EXPECT().SaveImage(gomock.Any(), gomock.Any()).Return("", errors.New("failed to save image"))
			},
			expectedURL: "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setupMocks()

			result, err := service.AddEvent(context.Background(), tc.event, tc.media)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedURL, result.ImageURL)
			}
		})
	}
}

func TestEventService_DeleteEvent(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventDB := mocks.NewMockEventDB(ctrl)
	mockImageDB := mocks.NewMockImageDB(ctrl)
	service := NewService(mockEventDB, mockImageDB)

	testCases := []struct {
		name        string
		ID          int
		AuthorID    int
		setupMocks  func()
		expectError bool
	}{
		{
			name:     "успешное выполнение",
			ID:       1,
			AuthorID: 1,
			setupMocks: func() {
				mockEventDB.EXPECT().GetEventByID(gomock.Any(), 1).Return(models.Event{ID: 1, AuthorID: 1, ImageURL: "path/to/image.jpg"}, nil)
				mockImageDB.EXPECT().DeleteImage(gomock.Any(), "path/to/image.jpg").Return(nil)
				mockEventDB.EXPECT().DeleteEvent(gomock.Any(), 1).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "нет доступа",
			ID:       1,
			AuthorID: 2,
			setupMocks: func() {
				mockEventDB.EXPECT().GetEventByID(gomock.Any(), 1).Return(models.Event{ID: 1, AuthorID: 1}, nil)
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setupMocks()

			err := service.DeleteEvent(context.Background(), tc.ID, tc.AuthorID)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEventService_SearchEvents(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventDB := mocks.NewMockEventDB(ctrl)
	service := NewService(mockEventDB, nil)

	testCases := []struct {
		name         string
		searchParams models.SearchParams
		setupMocks   func()
		expectError  bool
	}{
		{
			name: "успешное выполнение",
			searchParams: models.SearchParams{
				Tags: []string{"Music", "Live"},
			},
			setupMocks: func() {
				mockEventDB.EXPECT().SearchEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Event{}, nil)
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setupMocks()

			_, err := service.SearchEvents(context.Background(), tc.searchParams, models.PaginationParams{})
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEventService_UpdateEvent(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventDB := mocks.NewMockEventDB(ctrl)
	mockImageDB := mocks.NewMockImageDB(ctrl)
	service := NewService(mockEventDB, mockImageDB)

	mockFile := io.NopCloser(bytes.NewReader([]byte("updated image data")))

	testCases := []struct {
		name        string
		event       models.Event
		media       models.MediaFile
		setupMocks  func()
		expectedURL string
		expectError bool
	}{
		{
			name: "успешное выполенение",
			event: models.Event{
				ID:       1,
				AuthorID: 1,
				Title:    "Updated Event",
			},
			media: models.MediaFile{
				Filename: "new_image.jpg",
				File:     mockFile,
			},
			setupMocks: func() {
				mockEventDB.EXPECT().GetEventByID(gomock.Any(), 1).Return(models.Event{ID: 1, AuthorID: 1, ImageURL: "old_image.jpg"}, nil)
				mockImageDB.EXPECT().SaveImage(gomock.Any(), gomock.Any()).Return("path/to/new_image.jpg", nil)
				mockEventDB.EXPECT().UpdateEvent(gomock.Any(), gomock.Any()).Return(models.Event{ID: 1, AuthorID: 1, ImageURL: "path/to/new_image.jpg"}, nil)
				mockImageDB.EXPECT().DeleteImage(gomock.Any(), "old_image.jpg").Return(nil)
			},
			expectedURL: "path/to/new_image.jpg",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setupMocks()

			result, err := service.UpdateEvent(context.Background(), tc.event, tc.media)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedURL, result.ImageURL)
			}
		})
	}
}
