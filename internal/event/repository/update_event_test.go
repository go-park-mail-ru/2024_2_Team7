package eventRepository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kudago/internal/models"
)

func TestEventRepository_UpdateEvent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	eventStart := time.Now().Add(10 * time.Hour)
	eventFinish := eventStart.Add(5 * time.Hour)
	tests := []struct {
		name          string
		updatedEvent  models.Event
		mockSetup     func(m pgxmock.PgxConnIface)
		expectedEvent models.Event
		expectErr     bool
	}{
		{
			name: "Ошибка обновления события",
			updatedEvent: models.Event{
				ID:          1,
				Title:       "New Title",
				Description: "New Description",
				EventStart:  eventStart.Format(time.RFC3339),
				EventEnd:    eventFinish.Format(time.RFC3339),
				Location:    "New Location",
				Capacity:    100,
				CategoryID:  2,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`UPDATE event SET`).
					WithArgs(1, "New Title", "New Description", anyTime{}, anyTime{}, "New Location", 100, 2, anyTime{}, 0.0, 0.0).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedEvent: models.Event{},
			expectErr:     true,
		},
		{
			name: "Событие не найдено для обновления",
			updatedEvent: models.Event{
				ID:          999,
				Title:       "Nonexistent Event",
				Description: "No event found",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`UPDATE event SET`).
					WithArgs(999, "Nonexistent Event", "No event found", anyTime{}, anyTime{}, "", 0, 0, anyTime{}, 0.0, 0.0).
					WillReturnRows(pgxmock.NewRows([]string{"id"})) // No rows returned
			},
			expectedEvent: models.Event{},
			expectErr:     true,
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

			updatedEvent, err := db.UpdateEvent(ctx, tt.updatedEvent)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedEvent, updatedEvent)
			}
		})
	}
}

// Helper for any time comparison in mock rows.
type anyTime struct{}

func (anyTime) Scan(src interface{}) error {
	return nil
}
