package authService

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"kudago/internal/models"
	"kudago/internal/service/auth/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFile := io.NopCloser(bytes.NewReader([]byte("mock image data")))
	mockUserDB := mocks.NewMockUserDB(ctrl)
	mockSessionDB := mocks.NewMockSessionDB(ctrl)
	mockImageDB := mocks.NewMockImageDB(ctrl)
	service := NewService(mockUserDB, mockSessionDB, mockImageDB)

	testCases := []struct {
		name         string
		data         models.NewUserData
		setupMocks   func()
		expectedUser models.User
		expectError  bool
	}{
		{
			name: "успешное выполнение",
			data: models.NewUserData{
				User: models.User{
					Username: "user1",
					Email:    "user1@example.com",
				},
				Media: models.MediaFile{
					Filename: "profile.jpg",
					File:     mockFile,
				},
			},
			setupMocks: func() {
				mockImageDB.EXPECT().SaveImage(gomock.Any(), gomock.Any()).Return("path/to/profile.jpg", nil)
				mockUserDB.EXPECT().UserExists(gomock.Any(), "user1", "user1@example.com").Return(false, nil)
				mockUserDB.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(models.User{
					ID:       1,
					Username: "user1",
					Email:    "user1@example.com",
					ImageURL: "path/to/profile.jpg",
				}, nil)
			},
			expectedUser: models.User{
				ID:       1,
				Username: "user1",
				Email:    "user1@example.com",
				ImageURL: "path/to/profile.jpg",
			},
			expectError: false,
		},
		{
			name: "ошибка при сохранении изображения",
			data: models.NewUserData{
				User: models.User{
					Username: "user2",
					Email:    "user2@example.com",
				},
				Media: models.MediaFile{
					Filename: "profile.jpg",
					File:     mockFile,
				},
			},
			setupMocks: func() {
				mockImageDB.EXPECT().SaveImage(gomock.Any(), gomock.Any()).Return("", errors.New("failed to save image"))
			},
			expectedUser: models.User{},
			expectError:  true,
		},
		{
			name: "повторяется почта",
			data: models.NewUserData{
				User: models.User{
					Username: "user3",
					Email:    "user3@example.com",
				},
			},
			setupMocks: func() {
				mockUserDB.EXPECT().UserExists(gomock.Any(), "user3", "user3@example.com").Return(true, nil)
			},
			expectedUser: models.User{},
			expectError:  true,
		},
		{
			name: "сохранение без картинки",
			data: models.NewUserData{
				User: models.User{
					Username: "user4",
					Email:    "user4@example.com",
				},
			},
			setupMocks: func() {
				mockUserDB.EXPECT().UserExists(gomock.Any(), "user4", "user4@example.com").Return(false, nil)
				mockUserDB.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(models.User{
					ID:       2,
					Username: "user4",
					Email:    "user4@example.com",
				}, nil)
			},
			expectedUser: models.User{
				ID:       2,
				Username: "user4",
				Email:    "user4@example.com",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setupMocks()

			user, err := service.Register(context.Background(), tc.data)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestAuthService_UpdateUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFile := io.NopCloser(bytes.NewReader([]byte("updated image data")))
	mockUserDB := mocks.NewMockUserDB(ctrl)
	mockImageDB := mocks.NewMockImageDB(ctrl)
	service := NewService(mockUserDB, nil, mockImageDB)

	testCases := []struct {
		name         string
		data         models.NewUserData
		setupMocks   func()
		expectedUser models.User
		expectError  bool
	}{
		{
			name: "успешное выполнение",
			data: models.NewUserData{
				User: models.User{
					ID:       1,
					Username: "newuser",
					Email:    "newuser@example.com",
				},
				Media: models.MediaFile{
					Filename: "new_profile.jpg",
					File:     mockFile,
				},
			},
			setupMocks: func() {
				mockUserDB.EXPECT().GetUserByID(gomock.Any(), 1).Return(models.User{
					ID:       1,
					Username: "user1",
					Email:    "user1@example.com",
					ImageURL: "old/path/to/image.jpg",
				}, nil)
				mockUserDB.EXPECT().CheckUsername(gomock.Any(), "newuser", 1).Return(false, nil)
				mockUserDB.EXPECT().CheckEmail(gomock.Any(), "newuser@example.com", 1).Return(false, nil)
				mockImageDB.EXPECT().SaveImage(gomock.Any(), gomock.Any()).Return("new/path/to/image.jpg", nil)
				mockUserDB.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(models.User{
					ID:       1,
					Username: "newuser",
					Email:    "newuser@example.com",
					ImageURL: "new/path/to/image.jpg",
				}, nil)
				mockImageDB.EXPECT().DeleteImage(gomock.Any(), "old/path/to/image.jpg").Return(nil)
			},
			expectedUser: models.User{
				ID:       1,
				Username: "newuser",
				Email:    "newuser@example.com",
				ImageURL: "new/path/to/image.jpg",
			},
			expectError: false,
		},
		{
			name: "повторяется имя пользователя",
			data: models.NewUserData{
				User: models.User{
					ID:       1,
					Username: "existinguser",
				},
			},
			setupMocks: func() {
				mockUserDB.EXPECT().GetUserByID(gomock.Any(), 1).Return(models.User{
					ID:       1,
					Username: "user1",
					Email:    "user1@example.com",
				}, nil)
				mockUserDB.EXPECT().CheckUsername(gomock.Any(), "existinguser", 1).Return(true, nil)
			},
			expectedUser: models.User{},
			expectError:  true,
		},
		{
			name: "успешное выполнение без картинки",
			data: models.NewUserData{
				User: models.User{
					ID:       2,
					Username: "newuser2",
					Email:    "newuser2@example.com",
				},
			},
			setupMocks: func() {
				mockUserDB.EXPECT().GetUserByID(gomock.Any(), 2).Return(models.User{
					ID:       2,
					Username: "user2",
					Email:    "user2@example.com",
				}, nil)
				mockUserDB.EXPECT().CheckUsername(gomock.Any(), "newuser2", 2).Return(false, nil)
				mockUserDB.EXPECT().CheckEmail(gomock.Any(), "newuser2@example.com", 2).Return(false, nil)
				mockUserDB.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(models.User{
					ID:       2,
					Username: "newuser2",
					Email:    "newuser2@example.com",
				}, nil)
			},
			expectedUser: models.User{
				ID:       2,
				Username: "newuser2",
				Email:    "newuser2@example.com",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setupMocks()

			user, err := service.UpdateUser(context.Background(), tc.data)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserDB := mocks.NewMockUserDB(ctrl)
	service := NewService(mockUserDB, nil, nil)

	mockUser := models.User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	mockUserDB.EXPECT().GetUserByID(gomock.Any(), 1).Return(mockUser, nil)

	user, err := service.GetUserByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
}

func TestAuthService_CheckCredentials(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserDB := mocks.NewMockUserDB(ctrl)
	service := NewService(mockUserDB, nil, nil)

	mockUser := models.User{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	mockUserDB.EXPECT().CheckCredentials(gomock.Any(), "testuser", "password").Return(mockUser, nil)

	user, err := service.CheckCredentials(context.Background(), models.Credentials{
		Username: "testuser",
		Password: "password",
	})
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
}
