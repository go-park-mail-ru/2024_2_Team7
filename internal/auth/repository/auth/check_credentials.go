package userRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const checkCredentialsQuery = `
	SELECT id, username, email, created_at, url_to_avatar
	FROM "USER"
	WHERE username = $1 AND password_hash = $2`

func (d UserDB) CheckCredentials(ctx context.Context, username, password string) (models.User, error) {
	var userInfo UserInfo
	err := d.Pool.QueryRow(ctx, checkCredentialsQuery, username, password).Scan(
		&userInfo.ID,
		&userInfo.Username,
		&userInfo.Email,
		&userInfo.CreatedAt,
		&userInfo.ImageURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", models.LevelDB, models.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	user := ToDomainUser(userInfo)
	return user, nil
}
