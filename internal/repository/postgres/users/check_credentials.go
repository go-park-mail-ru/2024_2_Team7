package userRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const checkCredentialsQuery = `
	SELECT id, username, email, created_at, url_to_avatar, password_hash
	FROM "USER"
	WHERE username = $1`

func (d UserDB) CheckCredentials(ctx context.Context, username, password string) (models.User, error) {
	var userInfo UserInfo
	var storedPassHash []byte

	err := d.pool.QueryRow(ctx, checkCredentialsQuery, username).Scan(
		&userInfo.ID,
		&userInfo.Username,
		&userInfo.Email,
		&userInfo.CreatedAt,
		&userInfo.ImageURL,
		&storedPassHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", models.LevelDB, models.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	if !checkPass(storedPassHash, password) {
		return models.User{}, fmt.Errorf("%s: %w", models.LevelDB, models.ErrInvalidPassword)
	}

	user := toDomainUser(userInfo)
	return user, nil
}
