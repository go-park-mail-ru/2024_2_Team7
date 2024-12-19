package userRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const getUserByEmailOrUsernameQuery = `
		SELECT id 
		FROM "USER" 
		WHERE (username = $1 OR email = $2)`

func (d *UserDB) UserExists(ctx context.Context, user models.User) (bool, error) {
	var exists int
	query := getUserByEmailOrUsernameQuery
	args := []interface{}{user.Username, user.Email}

	if user.ID > 0 {
		query += " AND id != $3"
		args = append(args, user.ID)
	}

	err := d.Pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	return true, nil
}
