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

func TestEventRepository_DeleteEvent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		eventID   int
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr bool
	}{
		{
			name:    "Успешное удаление",
			eventID: 1,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM event WHERE id=\$1`).
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			expectErr: false,
		},
		{
			name:    "Событие не найдено",
			eventID: 2,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM event WHERE id=\$1`).
					WithArgs(2).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			expectErr: false,
		},
		{
			name:    "ошибка",
			eventID: 3,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM event WHERE id=\$1`).
					WithArgs(3).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: true,
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

			err = db.DeleteEvent(context.Background(), tt.eventID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
