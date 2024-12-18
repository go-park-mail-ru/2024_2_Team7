package eventRepository

import (
	"context"
	"fmt"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"kudago/internal/models"
)

func TestEventRepository_CreateEvent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		event     models.Event
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr bool
	}{
		{
			name: "ошибочное создание",
			event: models.Event{
				Title:       "Test Event",
				Description: "A test event",
				EventStart:  "2024-01-01T10:00:00Z",
				EventEnd:    "2024-01-01T12:00:00Z",
				Location:    "Test Location",
				Capacity:    100,
				AuthorID:    1,
				CategoryID:  2,
				Latitude:    10.0,
				Longitude:   20.0,
				Tag:         []string{"tag1", "tag2"},
				ImageURL:    "http://example.com/image.jpg",
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectBegin()
				m.ExpectQuery(createEventQuery).
					WithArgs("Test Event", "A test event", "2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z", "Test Location", 100, 1, 2, 10.0, 20.0).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
				m.ExpectExec("INSERT INTO event_tag").
					WithArgs(1, "tag1").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				m.ExpectExec("INSERT INTO event_tag").
					WithArgs(1, "tag2").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				m.ExpectExec("INSERT INTO event_media").
					WithArgs(1, "http://example.com/image.jpg").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				m.ExpectCommit()
			},
			expectErr: true,
		},
		{
			name: "Ошибка при выполнении запроса",
			event: models.Event{
				Title:       "Test Event",
				Description: "A test event",
				EventStart:  "2024-01-01T10:00:00Z",
				EventEnd:    "2024-01-01T12:00:00Z",
				Location:    "Test Location",
				Capacity:    100,
				AuthorID:    1,
				CategoryID:  2,
				Latitude:    10.0,
				Longitude:   20.0,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectBegin()
				m.ExpectQuery(createEventQuery).
					WithArgs("Test Event", "A test event", "2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z", "Test Location", 100, 1, 2, 10.0, 20.0).
					WillReturnError(fmt.Errorf("database error"))
				m.ExpectRollback()
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

			createdEvent, err := db.CreateEvent(context.Background(), tt.event)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, createdEvent.ID)
			}
		})
	}
}
