package userRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const deleteSubscription = `
	DELETE FROM SUBSCRIPTION
	WHERE subscriber_id=$1 AND follows_id=$2`

func (db *UserDB) Unsubscribe(ctx context.Context, subscription models.Subscription) error {
	result, err := db.pool.Exec(ctx, deleteSubscription, subscription.SubscriberID, subscription.FollowsID)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}
