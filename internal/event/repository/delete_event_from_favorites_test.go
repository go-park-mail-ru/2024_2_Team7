package eventRepository

import (
	"context"
	"fmt"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kudago/internal/models"
)

func TestEventRepository_DeleteEventFromFavorites(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		event     models.FavoriteEvent
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr error
	}{
		{
			name: "Успешное удаление",
			event: models.FavoriteEvent{
				EventID: 1,
				UserID:  1,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM FAVORITE_EVENT WHERE user_id=\$1 AND event_id=\$2`).
					WithArgs(1, 1).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			expectErr: nil,
		},
		{
			name: "Событие не найдено",
			event: models.FavoriteEvent{
				EventID: 2,
				UserID:  1,
			}, mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM FAVORITE_EVENT WHERE user_id=\$1 AND event_id=\$2`).
					WithArgs(1, 2).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			expectErr: models.ErrNotFound,
		},
		{
			name: "ошибка",
			event: models.FavoriteEvent{
				EventID: 3,
				UserID:  1,
			}, mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM FAVORITE_EVENT WHERE user_id=$1 AND event_id=$2`).
					WithArgs(1, 3).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: fmt.Errorf("database error"),
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

			err = db.DeleteEventFromFavorites(context.Background(), tt.event)
			if tt.expectErr != nil {
				assert.Error(t, tt.expectErr, err)
			}
		})
	}
}
