package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kudago/internal/auth/repository/auth"
	"testing"

	"kudago/internal/models"
)

func TestUserDB_CheckCredentials(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name         string
		username     string
		password     string
		mockSetup    func(m pgxmock.PgxConnIface)
		expectedUser models.User
		expectErr    bool
	}{
		{
			name:     "Пользователь не найден",
			username: "testuser",
			password: "wrongpassword",
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id, username, email, created_at, url_to_avatar`).
					WithArgs("testuser", "wrongpassword").
					WillReturnError(pgx.ErrNoRows)
			},
			expectedUser: models.User{},
			expectErr:    true,
		},
		{
			name:     "Ошибка при выполнении запроса",
			username: "testuser",
			password: "password123",
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id, username, email, created_at, url_to_avatar`).
					WithArgs("testuser", "password123").
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedUser: models.User{},
			expectErr:    true,
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
			user, err := db.CheckCredentials(ctx, tt.username, tt.password)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}
