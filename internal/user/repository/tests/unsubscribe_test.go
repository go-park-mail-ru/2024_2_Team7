package tests

import (
	"context"
	"fmt"
	"testing"

	"kudago/internal/models"
	"kudago/internal/user/repository"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserDB_Unsubscribe(t *testing.T) {
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
			name: "Успешная отмена подписки",
			subscription: models.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM SUBSCRIPTION WHERE subscriber_id=\$1 AND follows_id=\$2`).
					WithArgs(1, 2).
					WillReturnResult(pgxmock.NewResult("DELETE", 1)) // 1 row affected
			},
			expectErr:     false,
			expectedError: nil,
		},
		{
			name: "Нет строк для удаления",
			subscription: models.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM SUBSCRIPTION WHERE subscriber_id=\$1 AND follows_id=\$2`).
					WithArgs(1, 2).
					WillReturnResult(pgxmock.NewResult("DELETE", 0)) // No rows affected
			},
			expectErr:     true,
			expectedError: models.ErrNotFound,
		},
		{
			name: "Ошибка при выполнении запроса",
			subscription: models.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM SUBSCRIPTION WHERE subscriber_id=\$1 AND follows_id=\$2`).
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

			// Call the Unsubscribe method
			err = db.Unsubscribe(ctx, tt.subscription)

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
