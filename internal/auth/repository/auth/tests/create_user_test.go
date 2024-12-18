package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"kudago/internal/auth/repository/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kudago/internal/models"
)

func TestUserDB_CreateUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name         string
		user         models.User
		mockSetup    func(m pgxmock.PgxConnIface)
		expectedUser models.User
		expectErr    bool
	}{
		{
			name: "Успешное создание пользователя",
			user: models.User{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
				ImageURL: "http://example.com/avatar.jpg",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`INSERT INTO "USER" \(username, email, password_hash, url_to_avatar\)`).
					WithArgs("newuser", "newuser@example.com", "password123", "http://example.com/avatar.jpg").
					WillReturnRows(pgxmock.NewRows([]string{"id", "created_at"}).
						AddRow(1, time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)))
			},
			expectedUser: models.User{
				ID:       1,
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "",
				ImageURL: "http://example.com/avatar.jpg",
			},
			expectErr: false,
		},
		{
			name: "Ошибка при создании пользователя",
			user: models.User{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
				ImageURL: "http://example.com/avatar.jpg",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				// Мокируем ошибку при выполнении запроса
				m.ExpectQuery(`INSERT INTO "USER" \(username, email, password_hash, url_to_avatar\)`).
					WithArgs("newuser", "newuser@example.com", "password123", "http://example.com/avatar.jpg").
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedUser: models.User{},
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Создание мокированного соединения
			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			// Настройка моков
			tt.mockSetup(mockConn)

			// Инициализация репозитория с мокированным соединением
			db := userRepository.UserDB{Pool: mockConn}

			// Вызов функции
			user, err := db.CreateUser(ctx, tt.user)

			// Проверка ошибок
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}
