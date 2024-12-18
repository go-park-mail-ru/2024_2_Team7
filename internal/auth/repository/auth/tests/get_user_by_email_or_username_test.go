package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"kudago/internal/auth/repository/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kudago/internal/models"
)

func TestUserDB_UserExists(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		user      models.User
		mockSetup func(m pgxmock.PgxConnIface)
		expected  bool
		expectErr bool
	}{
		{
			name: "Пользователь существует по username",
			user: models.User{
				Username: "existinguser",
				Email:    "existinguser@example.com",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id FROM "USER" WHERE \(username = \$1 OR email = \$2\)`).
					WithArgs("existinguser", "existinguser@example.com").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
			},
			expected:  true,
			expectErr: false,
		},
		{
			name: "Пользователь не существует",
			user: models.User{
				Username: "nonexistinguser",
				Email:    "nonexistinguser@example.com",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id FROM "USER" WHERE \(username = \$1 OR email = \$2\)`).
					WithArgs("nonexistinguser", "nonexistinguser@example.com").
					WillReturnError(pgx.ErrNoRows)
			},
			expected:  false,
			expectErr: false,
		},
		{
			name: "Ошибка при выполнении запроса",
			user: models.User{
				Username: "user_with_error",
				Email:    "user_with_error@example.com",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id FROM "USER" WHERE \(username = \$1 OR email = \$2\)`).
					WithArgs("user_with_error", "user_with_error@example.com").
					WillReturnError(fmt.Errorf("database error"))
			},
			expected:  false,
			expectErr: true,
		},
		{
			name: "Пользователь существует, но исключаем его по ID",
			user: models.User{
				ID:       2,
				Username: "existinguser",
				Email:    "existinguser@example.com",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id FROM "USER" WHERE \(username = \$1 OR email = \$2\) AND id != \$3`).
					WithArgs("existinguser", "existinguser@example.com", 2).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
			},
			expected:  true,
			expectErr: false,
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

			exists, err := db.UserExists(ctx, tt.user)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			}
		})
	}
}
