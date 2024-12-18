package userRepository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kudago/internal/models"
)

func TestUserRepository_GetUserByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		userID     int
		mockSetup  func(m pgxmock.PgxConnIface)
		expectUser models.User
		expectErr  error
	}{
		{
			name:   "Успешное получение пользователя",
			userID: 1,
			mockSetup: func(m pgxmock.PgxConnIface) {
				rows := pgxmock.NewRows([]string{"id", "username", "email", "url_to_avatar"}).
					AddRow(1, "test_user", "test@example.com", "http://example.com/avatar.png")
				m.ExpectQuery(`SELECT id, username, email, url_to_avatar FROM "USER" WHERE id=\$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectUser: models.User{
				ID:       1,
				Username: "test_user",
				Email:    "test@example.com",
			},
			expectErr: nil,
		},
		{
			name:   "Пользователь не найден",
			userID: 2,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id, username, email, url_to_avatar FROM "USER" WHERE id=\$1`).
					WithArgs(2).
					WillReturnError(pgx.ErrNoRows)
			},
			expectUser: models.User{},
			expectErr:  models.ErrUserNotFound,
		},
		{
			name:   "Ошибка базы данных",
			userID: 3,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id, username, email, url_to_avatar FROM "USER" WHERE id=\$1`).
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			expectUser: models.User{},
			expectErr:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			db := NewDB(mockConn)

			user, err := db.GetUserByID(ctx, tt.userID)

			assert.Equal(t, tt.expectUser, user)
			assert.Error(t, err, tt.expectErr)
		})
	}
}
