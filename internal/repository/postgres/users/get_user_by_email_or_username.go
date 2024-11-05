package userRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const getUserByEmailOrUsernameQuery = `SELECT 1 FROM "USER" WHERE email=$1 OR username = $2 LIMIT 1`

func (d *UserDB) UserExists(ctx context.Context, username, email string) (bool, error) {
	var exists int
	err := d.pool.QueryRow(ctx, getUserByEmailOrUsernameQuery, email, username).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	return true, nil
}
