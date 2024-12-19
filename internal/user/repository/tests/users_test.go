package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"kudago/internal/models"
	"kudago/internal/user/repository"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_nilIfEmpty(t *testing.T) {
	t.Parallel()

	// Проверка с пустой строкой
	result := userRepository.NilIfEmpty("")
	assert.Nil(t, result)

	// Проверка с непустой строкой
	value := "non-empty"
	result = userRepository.NilIfEmpty(value)
	assert.NotNil(t, result)
	assert.Equal(t, value, *result)
}

func TestUserRepository_toDomainUser(t *testing.T) {
	t.Parallel()

	// Пример данных UserInfo с пустым ImageURL
	userInfo := userRepository.UserInfo{
		ID:         1,
		Username:   "testUser",
		Email:      "test@example.com",
		ImageURL:   nil,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	// Преобразование в models.User
	user := userRepository.ToDomainUser(userInfo)

	// Проверка значений
	assert.Equal(t, userInfo.ID, user.ID)
	assert.Equal(t, userInfo.Username, user.Username)
	assert.Equal(t, userInfo.Email, user.Email)
	assert.Equal(t, "", user.ImageURL) // Пустой ImageURL, так как он был nil

	// Пример данных UserInfo с непустым ImageURL
	imageURL := "http://example.com/avatar.jpg"
	userInfoWithImage := userRepository.UserInfo{
		ID:         2,
		Username:   "testUser2",
		Email:      "test2@example.com",
		ImageURL:   &imageURL,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	// Преобразование в models.User
	userWithImage := userRepository.ToDomainUser(userInfoWithImage)

	// Проверка значений
	assert.Equal(t, userInfoWithImage.ID, userWithImage.ID)
	assert.Equal(t, userInfoWithImage.Username, userWithImage.Username)
	assert.Equal(t, userInfoWithImage.Email, userWithImage.Email)
	assert.Equal(t, *userInfoWithImage.ImageURL, userWithImage.ImageURL) // Проверка, что ImageURL сохранен
}

func TestUserRepository_UpdateUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		updatedUser models.User
		mockSetup   func(m pgxmock.PgxConnIface)
		expectErr   bool
	}{
		{
			name: "Ошибка при обновлении",
			updatedUser: models.User{
				ID:       3,
				Username: "invalidUser",
				Email:    "invalid@example.com",
				ImageURL: "http://example.com/avatar.jpg",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`UPDATE "USER" SET username = COALESCE\(\$2, username\), email = COALESCE\(\$3, email\), URL_to_avatar = COALESCE\(\$4, URL_to_avatar\), modified_at = NOW\(\) WHERE id = \$1 RETURNING id, username, email, URL_to_avatar`).
					WithArgs(3, "invalidUser", "invalid@example.com", "http://example.com/avatar.jpg").
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			db := userRepository.UserDB{Pool: mockConn}

			user, err := db.UpdateUser(ctx, tt.updatedUser)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.updatedUser.Username, user.Username)
				assert.Equal(t, tt.updatedUser.Email, user.Email)
				assert.Equal(t, tt.updatedUser.ImageURL, user.ImageURL)
			}
		})
	}
}
