package tests

import (
	"context"
	"fmt"
	"testing"

	"kudago/internal/models"
	"kudago/internal/user/repository"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserDB_GetSubscriptions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		ID        int
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr bool
		expectRes []models.User
	}{
		{
			name: "Успешное получение подписок",
			ID:   1,
			mockSetup: func(m pgxmock.PgxConnIface) {
				rows := pgxmock.NewRows([]string{"id", "username", "email", "url_to_avatar"}).
					AddRow(1, "user1", "user1@example.com", "https://avatar.com/user1").
					AddRow(2, "user2", "user2@example.com", "https://avatar.com/user2")

				m.ExpectQuery(`SELECT u.id, u.username, u.email, u.url_to_avatar FROM "USER" u JOIN SUBSCRIPTION s ON s.follows_id = u.id WHERE s.subscriber_id = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectErr: false,
			expectRes: []models.User{
				{ID: 1, Username: "user1", Email: "user1@example.com", ImageURL: "https://avatar.com/user1"},
				{ID: 2, Username: "user2", Email: "user2@example.com", ImageURL: "https://avatar.com/user2"},
			},
		},
		{
			name: "Ошибка при запросе",
			ID:   3,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT u.id, u.username, u.email, u.url_to_avatar FROM "USER" u JOIN SUBSCRIPTION s ON s.follows_id = u.id WHERE s.subscriber_id = \$1`).
					WithArgs(3).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: true,
			expectRes: nil,
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

			res, err := db.GetSubscriptions(ctx, tt.ID)

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
