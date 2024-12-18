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

func TestUserRepository_GetUserByEmailOrUsername(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		user       models.User
		mockSetup  func(m pgxmock.PgxConnIface)
		expectUser bool
		expectErr  error
	}{
		{
			name: "Успешное получение пользователя",
			user: models.User{
				Username: "user1",
				Email:    "user1",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(1)
				m.ExpectQuery(`SELECT id\s+FROM "USER"\s+WHERE \(username = \$1 OR email = \$2\)`).
					WithArgs("user1", "user1").
					WillReturnRows(rows)
			},
			expectUser: true,
			expectErr:  nil,
		},
		{
			name: "Пользователь не найден",
			user: models.User{
				Username: "user1",
				Email:    "user1",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id\s+FROM "USER"\s+WHERE \(username = \$1 OR email = \$2\)`).
					WithArgs("user1", "user1").
					WillReturnError(pgx.ErrNoRows)
			},
			expectUser: false,
			expectErr:  nil,
		},
		{
			name: "Ошибка базы данных",
			user: models.User{
				Username: "user1",
				Email:    "user1",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id FROM "USER" WHERE (username = \$1 OR email = \$2)`).
					WithArgs("user1", "user1").
					WillReturnError(errors.New("database error"))
			},
			expectUser: false,
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

			user, err := db.UserExists(ctx, tt.user)

			assert.Equal(t, tt.expectUser, user)
			if tt.expectErr != nil {
				assert.Error(t, tt.expectErr, err)
			}
		})
	}
}
