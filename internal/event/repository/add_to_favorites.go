package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

const insertNewFavorite = `
	INSERT INTO FAVORITE_EVENT (user_id, event_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING`

func (db *EventDB) AddEventToFavorites(ctx context.Context, newFavorite models.FavoriteEvent) error {
	result, err := db.pool.Exec(ctx, insertNewFavorite, newFavorite.UserID, newFavorite.EventID)
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
