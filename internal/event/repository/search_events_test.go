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

func TestEventRepository_SearchEvents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name             string
		params           models.SearchParams
		paginationParams models.PaginationParams
		mockSetup        func(m pgxmock.PgxConnIface)
		expectedEvents   []models.Event
		expectErr        bool
	}{
		{
			name: "поиск событий",
			params: models.SearchParams{
				Query:      "test",
				Category:   1,
				EventStart: time.Now().Format("2006-01-02 15:04:05"),                     // Convert time to string
				EventEnd:   time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"), // Convert time to string
				Tags:       []string{"tag1", "tag2"},
			},
			paginationParams: models.PaginationParams{
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
					event.location, event.capacity, event.created_at, event.user_id, event.category_id, event.lat, event.lon,
					COALESCE\(array_agg\(DISTINCT tag.name\) FILTER \(WHERE tag.name IS NOT NULL\), ARRAY\[\]::TEXT\[\]\) AS tags,
					COALESCE\(media_url.url, ''\) AS media_link	
				FROM event
				LEFT JOIN event_tag ON event.id = event_tag.event_id
				LEFT JOIN tag ON tag.id = event_tag.tag_id
				LEFT JOIN media_url ON event.id = media_url.event_id
				WHERE
					\(\$1::TEXT IS NULL OR event.title ILIKE '%' || \$1 || '%' OR event.description ILIKE '%' || \$1 || '%'\)
					AND \(\$2::INT IS NULL OR event.category_id = \$2\)
					AND \(\$3::TIMESTAMP IS NULL OR event.event_start >= \$3\)
					AND \(\$4::TIMESTAMP IS NULL OR event.event_finish <= \$4\)
					AND \(\$8::DOUBLE PRECISION IS NULL OR event.lat >= \$8\) -- Минимальная широта
					AND \(\$9::DOUBLE PRECISION IS NULL OR event.lat <= \$9\) -- Максимальная широта
					AND \(\$10::DOUBLE PRECISION IS NULL OR event.lon >= \$10\) -- Минимальная долгота
					AND \(\$11::DOUBLE PRECISION IS NULL OR event.lon <= \$11\) -- Максимальная долгота
				GROUP BY event.id, media_url.url
				HAVING (
					\$5::TEXT\[\] IS NULL 
					OR array_length\(\$5::TEXT\[\], 1\) = 0 
					OR array_length\(array_agg\(DISTINCT LOWER\(tag.name\)\), 1\) = 0 
					OR array_agg\(DISTINCT LOWER\(tag.name\)\) @> \$5::TEXT\[\]
				)
				ORDER BY event.event_finish ASC
				LIMIT \$6 OFFSET \$7;`).
					WithArgs("test", 1, time.Now().Format("2006-01-02 15:04:05"), time.Now().Add(24*time.Hour).Format("2006-01-02 15:04:05"), []string{"tag1", "tag2"}, 10, 0, nil, nil, nil, nil).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "title", "description", "event_start", "event_finish", "location", "capacity", "created_at", "user_id", "category_id", "lat", "lon", "tags", "media_link",
					}).
						AddRow(1, "Event 1", "Description 1", time.Now().Format("2006-01-02 15:04:05"), time.Now().Add(1*time.Hour).Format("2006-01-02 15:04:05"), "Location 1", 100, time.Now(), 1, 1, 10.0, 20.0, []string{"tag1", "tag2"}, "http://example.com").
						AddRow(2, "Event 2", "Description 2", time.Now().Format("2006-01-02 15:04:05"), time.Now().Add(2*time.Hour).Format("2006-01-02 15:04:05"), "Location 2", 200, time.Now(), 2, 2, 15.0, 25.0, []string{"tag2"}, "http://example2.com"))
			},
			expectedEvents: []models.Event{
				{
					ID:          1,
					Title:       "Event 1",
					Description: "Description 1",
					EventStart:  time.Now().Format("2006-01-02 15:04:05"),
					EventEnd:    time.Now().Add(1 * time.Hour).Format("2006-01-02 15:04:05"),
					Location:    "Location 1",
					Capacity:    100,
					CreatedAt:   time.Now(),
					CategoryID:  1,
					AuthorID:    1,
					Latitude:    10.0,
					Longitude:   20.0,
					Tag:         []string{"tag1", "tag2"},
					ImageURL:    "http://example.com",
				},
				{
					ID:          2,
					Title:       "Event 2",
					Description: "Description 2",
					EventStart:  time.Now().Format("2006-01-02 15:04:05"),
					EventEnd:    time.Now().Add(2 * time.Hour).Format("2006-01-02 15:04:05"),
					Location:    "Location 2",
					Capacity:    200,
					CreatedAt:   time.Now(),
					CategoryID:  2,
					AuthorID:    2,
					Latitude:    15.0,
					Longitude:   25.0,
					Tag:         []string{"tag2"},
					ImageURL:    "http://example2.com",
				},
			},
			expectErr: true,
		},
		{
			name: "Ошибка при поиске",
			params: models.SearchParams{
				Query: "test",
			},
			paginationParams: models.PaginationParams{
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
					event.location, event.capacity, event.created_at, event.user_id, event.category_id, event.lat, event.lon,
					COALESCE\(array_agg\(DISTINCT tag.name\) FILTER \(WHERE tag.name IS NOT NULL\), ARRAY\[\]::TEXT\[\]\) AS tags,
					COALESCE\(media_url.url, ''\) AS media_link	
				FROM event
				LEFT JOIN event_tag ON event.id = event_tag.event_id
				LEFT JOIN tag ON tag.id = event_tag.tag_id
				LEFT JOIN media_url ON event.id = media_url.event_id
				WHERE
					\(\$1::TEXT IS NULL OR event.title ILIKE '%' || \$1 || '%' OR event.description ILIKE '%' || \$1 || '%'\)
					AND \(\$2::INT IS NULL OR event.category_id = \$2\)
					AND \(\$3::TIMESTAMP IS NULL OR event.event_start >= \$3\)
					AND \(\$4::TIMESTAMP IS NULL OR event.event_finish <= \$4\)
					AND \(\$8::DOUBLE PRECISION IS NULL OR event.lat >= \$8\) -- Минимальная широта
					AND \(\$9::DOUBLE PRECISION IS NULL OR event.lat <= \$9\) -- Максимальная широта
					AND \(\$10::DOUBLE PRECISION IS NULL OR event.lon >= \$10\) -- Минимальная долгота
					AND \(\$11::DOUBLE PRECISION IS NULL OR event.lon <= \$11\) -- Максимальная долгота
				GROUP BY event.id, media_url.url
				HAVING (
					\$5::TEXT\[\] IS NULL 
					OR array_length\(\$5::TEXT\[\], 1\) = 0 
					OR array_length\(array_agg\(DISTINCT LOWER\(tag.name\)\), 1\) = 0 
					OR array_agg\(DISTINCT LOWER\(tag.name\)\) @> \$5::TEXT\[\]
				)
				ORDER BY event.event_finish ASC
				LIMIT \$6 OFFSET \$7;`).
					WithArgs("test", nil, nil, nil, nil, 10, 0, nil, nil, nil, nil).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedEvents: nil,
			expectErr:      true,
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

			events, err := db.SearchEvents(ctx, tt.params, tt.paginationParams)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedEvents, events)
			}
		})
	}
}
