package service

import (
	"context"
	"testing"

	"kudago/internal/event/service/mocks"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEventService_DeleteEvent(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventDB := mocks.NewMockEventDB(ctrl)
	service := NewService(mockEventDB)

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
	service := NewService(mockEventDB)

	testCases := []struct {
		name         string
		searchParams models.SearchParams
		setupMocks   func()
		expectError  bool
	}{
		{
			name: "success search",
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

// func TestEventService_UpdateEvent(t *testing.T) {
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockEventDB := mocks.NewMockEventDB(ctrl)
// 	service := NewService(mockEventDB)

// 	event := models.Event{
// 		ID:       1,
// 		AuthorID: 1,
// 		Title:    "Event",
// 	}

// 	testCases := []struct {
// 		name          string
// 		event         models.Event
// 		setupMocks    func()
// 		expectedEvent models.Event
// 		expectedError error
// 	}{
// 		{
// 			name:  "success update event",
// 			event: event,
// 			setupMocks: func() {
// 				mockEventDB.EXPECT().GetEventByID(gomock.Any(), 1).Return(event, nil)
// 				mockEventDB.EXPECT().UpdateEvent(gomock.Any(), gomock.Any()).Return(event, nil)
// 			},
// 			expectedEvent: event,
// 			expectedError: nil,
// 		},
// 		{
// 			name:  "event doesn't exists",
// 			event: event,
// 			setupMocks: func() {
// 				mockEventDB.EXPECT().GetEventByID(gomock.Any(), 1).Return(models.Event{}, models.ErrEventNotFound)
// 			},
// 			expectedEvent: models.Event{},
// 			expectedError: models.ErrEventNotFound,
// 		},
// 		{
// 			name: "access denied",
// 			event: models.Event{
// 				ID:       1,
// 				AuthorID: 3,
// 				Title:    "access denied",
// 			},
// 			setupMocks: func() {
// 				mockEventDB.EXPECT().GetEventByID(gomock.Any(), 1).Return(event, nil)
// 			},
// 			expectedEvent: models.Event{},
// 			expectedError: models.ErrAccessDenied,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.setupMocks()

// 			result, err := service.UpdateEvent(context.Background(), tc.event)
// 			assert.ErrorAs(t, tc.expectedError, err)
// 			assert.Equal(t, tc.expectedEvent, result)
// 		})
// 	}
// }
