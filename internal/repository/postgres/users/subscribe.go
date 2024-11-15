package userRepository

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"kudago/internal/models"

	"github.com/jackc/pgconn"
)

const insertSubscription = `
	INSERT INTO SUBSCRIPTION (subscriber_id, subscribed_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING`

func (db *UserDB) Subscribe(ctx context.Context, subscription models.Subscription) error {
	result, err := db.pool.Exec(ctx, insertSubscription, subscription.SubscriberID, subscription.FollowsID)
	if err != nil {
		var pgErr *pgconn.PgError
		fmt.Println(errors.As(err, &pgErr), reflect.TypeOf(err))
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
