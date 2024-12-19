package userRepository

import (
	"context"
	"errors"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

const insertSubscription = `
	INSERT INTO SUBSCRIPTION (subscriber_id, follows_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING`

func (db *UserDB) Subscribe(ctx context.Context, subscription models.Subscription) error {
	result, err := db.Pool.Exec(ctx, insertSubscription, subscription.SubscriberID, subscription.FollowsID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return models.ErrForeignKeyViolation
		}
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNothingToInsert
	}
	return nil
}
