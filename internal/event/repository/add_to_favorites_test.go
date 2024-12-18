package eventRepository

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"kudago/internal/models"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventDB_AddEventToFavorites(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name          string
		newFavorite   models.FavoriteEvent
		mockSetup     func(m pgxmock.PgxConnIface)
		expectedError error
	}{
		{
			name: "Успешное добавление",
			newFavorite: models.FavoriteEvent{
				UserID:  1,
				EventID: 100,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO FAVORITE_EVENT \(user_id, event_id\) VALUES \(\$1, \$2\) ON CONFLICT DO NOTHING`).
					WithArgs(1, 100).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedError: nil,
		},
		{
			name: "Конфликт записи (ничего не вставлено)",
			newFavorite: models.FavoriteEvent{
				UserID:  2,
				EventID: 200,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO FAVORITE_EVENT \(user_id, event_id\) VALUES \(\$1, \$2\) ON CONFLICT DO NOTHING`).
					WithArgs(2, 200).
					WillReturnResult(pgxmock.NewResult("INSERT", 0))
			},
			expectedError: models.ErrNothingToInsert,
		},
		{
			name: "Ошибка внешнего ключа",
			newFavorite: models.FavoriteEvent{
				UserID:  3,
				EventID: 300,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO FAVORITE_EVENT \(user_id, event_id\) VALUES \(\$1, \$2\) ON CONFLICT DO NOTHING`).
					WithArgs(3, 300).
					WillReturnError(&pgconn.PgError{Code: "23503"}) // Код ошибки внешнего ключа
			},
			expectedError: models.ErrForeignKeyViolation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			db := EventDB{pool: mockConn}

			err = db.AddEventToFavorites(ctx, tt.newFavorite)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
