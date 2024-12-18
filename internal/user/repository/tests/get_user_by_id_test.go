package tests

import (
	"context"
	"kudago/internal/models"
	"kudago/internal/user/repository"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserDB_GetUserByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		ID        int
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr bool
		expectRes models.User
	}{

		{
			name: "Пользователь не найден",
			ID:   2,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id, username, email, url_to_avatar FROM "USER" WHERE id=\$1`).
					WithArgs(2).
					WillReturnError(pgx.ErrNoRows)
			},
			expectErr: true,
			expectRes: models.User{},
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

			res, err := db.GetUserByID(ctx, tt.ID)

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
