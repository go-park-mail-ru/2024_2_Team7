package tests

import (
	"context"
	"fmt"
	"kudago/internal/models"
	"kudago/internal/user/repository"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserDB_UserExists(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		user      models.User
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr bool
		expectRes bool
	}{
		{
			name: "Пользователь существует по username",
			user: models.User{
				Username: "user1",
				Email:    "user1@example.com",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id FROM "USER" WHERE \(username = \$1 OR email = \$2\)`).
					WithArgs("user1", "user1@example.com").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectErr: false,
			expectRes: true,
		},
		{
			name: "Пользователь не существует",
			user: models.User{
				Username: "nonexistent",
				Email:    "nonexistent@example.com",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id FROM "USER" WHERE \(username = \$1 OR email = \$2\)`).
					WithArgs("nonexistent", "nonexistent@example.com").
					WillReturnError(pgx.ErrNoRows)
			},
			expectErr: false,
			expectRes: false,
		},
		{
			name: "Ошибка при выполнении запроса",
			user: models.User{
				Username: "user1",
				Email:    "user1@example.com",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id FROM "USER" WHERE \(username = \$1 OR email = \$2\)`).
					WithArgs("user1", "user1@example.com").
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: true,
			expectRes: false,
		},
		{
			name: "Пользователь существует с исключением ID",
			user: models.User{
				ID:       5,
				Username: "user1",
				Email:    "user1@example.com",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id FROM "USER" WHERE \(username = \$1 OR email = \$2\) AND id != \$3`).
					WithArgs("user1", "user1@example.com", 5).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectErr: false,
			expectRes: true,
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

			res, err := db.UserExists(ctx, tt.user)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectRes, res)
			}
		})
	}
}
