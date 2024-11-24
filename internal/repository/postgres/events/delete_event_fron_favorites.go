package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const deleteFavorite = `DELETE FROM FAVORITE_EVENT WHERE user_id=$1 AND event_id=$2`

func (db *EventDB) DeleteEventFromFavorites(ctx context.Context, favorite models.FavoriteEvent) error {
	_, err := db.pool.Exec(ctx, deleteFavorite, favorite.UserID, favorite.EventID)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	return nil
}
