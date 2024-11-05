package eventRepository

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// 	"kudago/internal/models"
// )

// func TestEventDB_AddEvent(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("unexpected error: %s", err)
// 	}
// 	defer db.Close()

// 	eventDB := &EventDB{pool: db}
// 	ctx := context.Background()

// 	event := models.Event{
// 		Title:       "Test Event",
// 		Description: "Event Description",
// 		EventStart:  time.Now().Format(time.RFC3339),
// 		EventEnd:    time.Now().Add(2 * time.Hour).Format(time.RFC3339),
// 		Location:    "Location A",
// 		Capacity:    100,
// 		AuthorID:    1,
// 		CategoryID:  1,
// 		Tag:         []string{"music", "outdoor"},
// 		ImageURL:    "image/path.jpg",
// 	}

// 	testCases := []struct {
// 		name        string
// 		setupMocks  func()
// 		expectError bool
// 	}{
// 		{
// 			name: "successful insertion",
// 			setupMocks: func() {
// 				mock.ExpectBegin()

// 				// Expect event insertion query
// 				mock.ExpectQuery(addEventQuery).
// 					WithArgs(event.Title, event.Description, event.EventStart, event.EventEnd, event.Location, event.Capacity, event.AuthorID, event.CategoryID).
// 					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

// 				// Expect tag insertion queries
// 				mock.ExpectExec(insertTagsQuery).WithArgs("music").WillReturnResult(sqlmock.NewResult(1, 1))
// 				mock.ExpectExec(insertTagsQuery).WithArgs("outdoor").WillReturnResult(sqlmock.NewResult(1, 1))
// 				mock.ExpectQuery(selectTagIDsQuery).WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
// 				mock.ExpectExec(insertEventTagQuery).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
// 				mock.ExpectExec(insertEventTagQuery).WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))

// 				// Expect media URL insertion
// 				mock.ExpectExec(insertMediaURL).WithArgs(1, event.ImageURL).WillReturnResult(sqlmock.NewResult(1, 1))

// 				mock.ExpectCommit()
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "error inserting tags",
// 			setupMocks: func() {
// 				mock.ExpectBegin()

// 				mock.ExpectQuery(addEventQuery).
// 					WithArgs(event.Title, event.Description, event.EventStart, event.EventEnd, event.Location, event.Capacity, event.AuthorID, event.CategoryID).
// 					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

// 				// Simulate failure when adding a tag
// 				mock.ExpectExec(insertTagsQuery).WithArgs("music").WillReturnError(errors.New("failed to insert tag"))

// 				mock.ExpectRollback()
// 			},
// 			expectError: true,
// 		},
// 		{
// 			name: "error inserting media URL",
// 			setupMocks: func() {
// 				mock.ExpectBegin()

// 				mock.ExpectQuery(addEventQuery).
// 					WithArgs(event.Title, event.Description, event.EventStart, event.EventEnd, event.Location, event.Capacity, event.AuthorID, event.CategoryID).
// 					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

// 				// Expect tag insertion to be successful
// 				mock.ExpectExec(insertTagsQuery).WithArgs("music").WillReturnResult(sqlmock.NewResult(1, 1))
// 				mock.ExpectExec(insertTagsQuery).WithArgs("outdoor").WillReturnResult(sqlmock.NewResult(1, 1))
// 				mock.ExpectQuery(selectTagIDsQuery).WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
// 				mock.ExpectExec(insertEventTagQuery).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
// 				mock.ExpectExec(insertEventTagQuery).WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))

// 				// Simulate failure when adding media URL
// 				mock.ExpectExec(insertMediaURL).WithArgs(1, event.ImageURL).WillReturnError(errors.New("failed to insert media URL"))

// 				mock.ExpectRollback()
// 			},
// 			expectError: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.setupMocks()

// 			result, err := eventDB.AddEvent(ctx, event)
// 			if tc.expectError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, 1, result.ID) // Confirm ID is set as expected
// 			}

// 			// Ensure all expectations were met
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("unfulfilled expectations: %s", err)
// 			}
// 		})
// 	}
// }
