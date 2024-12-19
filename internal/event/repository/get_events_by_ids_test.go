package eventRepository

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventRepository_GetEventsByIDs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		IDs       []int
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr error
	}{
		{
			name: "Успешное получение событий по ID",
			IDs:  []int{1, 2},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
					event.location, event.capacity, event.created_at, event.user_id, event.category_id, event.lat, event.lon, 
					COALESCE\(array_agg\(COALESCE\(tag.name, ''\), \{\}\)\) AS tags, media_url.url AS media_link`).
					WithArgs([]int{1, 2}).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "title", "description", "event_start", "event_finish", "location", "capacity",
						"created_at", "user_id", "category_id", "lat", "lon", "tags", "media_link",
					}).
						AddRow(1, "Event 1", "Description 1", "2024-01-01", "2024-01-02", "Location 1", 100,
							"2024-01-01", 1, 1, 55.7558, 37.6173, "{tag1, tag2}", "url1").
						AddRow(2, "Event 2", "Description 2", "2024-02-01", "2024-02-02", "Location 2", 200,
							"2024-02-01", 2, 2, 40.7128, -74.0060, "{tag3, tag4}", "url2"))
			},
			expectErr: nil,
		},
		{
			name: "События не найдены",
			IDs:  []int{999},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
					event.location, event.capacity, event.created_at, event.user_id, event.category_id, event.lat, event.lon, 
					COALESCE\(array_agg\(COALESCE\(tag.name, ''\), \{\}\)\) AS tags, media_url.url AS media_link`).
					WithArgs([]int{999}).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "title", "description", "event_start", "event_finish", "location", "capacity",
						"created_at", "user_id", "category_id", "lat", "lon", "tags", "media_link",
					}))
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
