package eventRepository

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"kudago/internal/models"

// 	"github.com/pashagolub/pgxmock/v4"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestEventsRepository_GetEventsByCategory(t *testing.T) {
// 	t.Parallel()

// 	ctx := context.Background()
// 	eventStart := time.Now().Add(10 * time.Hour)
// 	eventFinish := eventStart.Add(5 * time.Hour)

// 	tests := []struct {
// 		name           string
// 		categoryID     int
// 		pagination     models.PaginationParams
// 		mockSetup      func(m pgxmock.PgxConnIface)
// 		expectErr      bool
// 		expectedEvents []models.Event
// 	}{
// 		{
// 			name:       "успешное выполнение",
// 			categoryID: 2,
// 			pagination: models.PaginationParams{
// 				Limit:  2,
// 				Offset: 0,
// 			},
// 			mockSetup: func(m pgxmock.PgxConnIface) {
// 				img1 := "http://example.com/image1.jpg"
// 				img2 := "http://example.com/image2.jpg"
// 				rows := m.NewRows([]string{
// 					"id", "title", "description", "event_start", "event_finish",
// 					"location", "capacity", "created_at", "user_id", "category_id", "tags", "media_link",
// 				}).AddRow(
// 					1, "Music Festival", "A great music festival", eventStart, eventFinish,
// 					"Central Park", 1000, time.Now(), 1, 2, []string{"rock", "pop"}, &img1,
// 				).AddRow(
// 					2, "Art Exhibition", "Amazing art", eventStart, eventFinish,
// 					"Art Gallery", 200, time.Now(), 2, 3, []string{"art", "gallery"}, &img2,
// 				)
// 				m.ExpectQuery(`SELECT event.id, event.title, event.description, event.event_start, event.event_finish`).
// 					WithArgs(2, 2, 0).
// 					WillReturnRows(rows)
// 			},
// 			expectErr: false,
// 			expectedEvents: []models.Event{
// 				{
// 					ID:          1,
// 					Title:       "Music Festival",
// 					Description: "A great music festival",
// 					EventStart:  eventStart.Format(time.RFC3339),
// 					EventEnd:    eventFinish.Format(time.RFC3339),
// 					Location:    "Central Park",
// 					Capacity:    1000,
// 					AuthorID:    1,
// 					CategoryID:  2,
// 					Tag:         []string{"rock", "pop"},
// 					ImageURL:    "http://example.com/image1.jpg",
// 				},
// 				{
// 					ID:          2,
// 					Title:       "Art Exhibition",
// 					Description: "Amazing art",
// 					EventStart:  eventStart.Format(time.RFC3339),
// 					EventEnd:    eventFinish.Format(time.RFC3339),
// 					Location:    "Art Gallery",
// 					Capacity:    200,
// 					AuthorID:    2,
// 					CategoryID:  3,
// 					Tag:         []string{"art", "gallery"},
// 					ImageURL:    "http://example.com/image2.jpg",
// 				},
// 			},
// 		},
// 		{
// 			name:       "query execution error",
// 			categoryID: 2,
// 			pagination: models.PaginationParams{
// 				Limit:  2,
// 				Offset: 0,
// 			},
// 			mockSetup: func(m pgxmock.PgxConnIface) {
// 				m.ExpectQuery(`SELECT event.id, event.title, event.description, event.event_start, event.event_finish`).
// 					WithArgs(2, 2, 0).
// 					WillReturnError(errors.New("query error"))
// 			},
// 			expectErr:      true,
// 			expectedEvents: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			mockConn, err := pgxmock.NewConn()
// 			require.NoError(t, err)
// 			defer mockConn.Close(ctx)

// 			tt.mockSetup(mockConn)

// 			db := &EventDB{pool: mockConn}

// 			events, err := db.GetEventsByCategory(ctx, tt.categoryID, tt.pagination)

// 			if tt.expectErr {
// 				assert.Error(t, err)
// 				assert.Empty(t, events)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expectedEvents, events)
// 			}
// 		})
// 	}
// }
