package tests

import (
	"context"
	"fmt"
	"kudago/internal/models"
	"kudago/internal/user/repository"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserDB_Subscribe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name          string
		subscription  models.Subscription
		mockSetup     func(m pgxmock.PgxConnIface)
		expectErr     bool
		expectedError error
	}{
		{
			name: "Успешная подписка",
			subscription: models.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO SUBSCRIPTION \(subscriber_id, follows_id\) VALUES \(\$1, \$2\)`).
					WithArgs(1, 2).
					WillReturnResult(pgxmock.NewResult("INSERT", 1)) // 1 row inserted
			},
			expectErr:     false,
			expectedError: nil,
		},
		{
			name: "Нарушение внешнего ключа",
			subscription: models.Subscription{
				SubscriberID: 999, // Assuming this ID doesn't exist in the DB
				FollowsID:    2,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO SUBSCRIPTION \(subscriber_id, follows_id\) VALUES \(\$1, \$2\)`).
					WithArgs(999, 2).
					WillReturnError(&pgconn.PgError{Code: "23503"}) // Foreign key violation error
			},
			expectErr:     true,
			expectedError: models.ErrForeignKeyViolation,
		},
		{
			name: "Нет строк для вставки",
			subscription: models.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO SUBSCRIPTION \(subscriber_id, follows_id\) VALUES \(\$1, \$2\)`).
					WithArgs(1, 2).
					WillReturnResult(pgxmock.NewResult("INSERT", 0)) // No rows inserted
			},
			expectErr:     true,
			expectedError: models.ErrNothingToInsert,
		},
		{
			name: "Ошибка при выполнении запроса",
			subscription: models.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO SUBSCRIPTION \(subscriber_id, follows_id\) VALUES \(\$1, \$2\)`).
					WithArgs(1, 2).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr:     true,
			expectedError: fmt.Errorf("%s: %w", models.LevelDB, fmt.Errorf("database error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock connection
			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			// Setup mock behavior
			tt.mockSetup(mockConn)

			// Create the repository with the mock connection
			db := userRepository.UserDB{Pool: mockConn}

			// Call the Subscribe method
			err = db.Subscribe(ctx, tt.subscription)

			// Check for the expected error
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
