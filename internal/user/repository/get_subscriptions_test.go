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

func TestUserRepository_GetSubscriptions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		userID     int
		mockSetup  func(m pgxmock.PgxConnIface)
		expectUser []models.User
		expectErr  error
	}{
		{
			name:   "Успешное получение пользователя",
			userID: 1,
			mockSetup: func(m pgxmock.PgxConnIface) {
				rows := pgxmock.NewRows([]string{"id", "username", "email", "url_to_avatar"}).
					AddRow(1, "test_user", "test@example.com", "http://example.com/avatar.png")
				m.ExpectQuery(`SELECT\s+u\.id,\s+u\.username,\s+u\.email,\s+u\.url_to_avatar\s+FROM\s+"USER"\s+u\s+JOIN\s+SUBSCRIPTION\s+s\s+ON\s+s\.follows_id\s+=\s+u\.id\s+WHERE\s+s\.subscriber_id\s+=\s+\$1;`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectUser: []models.User{
				{
					ID:       1,
					Username: "test_user",
					Email:    "test@example.com",
					ImageURL: "http://example.com/avatar.png",
				},
			},
			expectErr: nil,
		},
		{
			name:   "Пользователь не найден",
			userID: 2,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT\s+u\.id,\s+u\.username,\s+u\.email,\s+u\.url_to_avatar\s+FROM\s+"USER"\s+u\s+JOIN\s+SUBSCRIPTION\s+s\s+ON\s+s\.follows_id\s+=\s+u\.id\s+WHERE\s+s\.subscriber_id\s+=\s+\$1;`).
					WithArgs(2).
					WillReturnError(pgx.ErrNoRows)
			},
			expectUser: nil,
			expectErr:  models.ErrUserNotFound,
		},
		{
			name:   "Ошибка базы данных",
			userID: 3,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT\s+u\.id,\s+u\.username,\s+u\.email,\s+u\.url_to_avatar\s+FROM\s+"USER"\s+u\s+JOIN\s+SUBSCRIPTION\s+s\s+ON\s+s\.follows_id\s+=\s+u\.id\s+WHERE\s+s\.subscriber_id\s+=\s+\$1;`).
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			expectUser: nil,
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

			user, err := db.GetSubscriptions(ctx, tt.userID)

			assert.Equal(t, tt.expectUser, user)
			if tt.expectErr != nil {
				assert.Error(t, err, tt.expectErr)
			}
		})
	}
}
