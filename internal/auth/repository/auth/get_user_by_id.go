package userRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
)

const getUserByIDQuery = `SELECT id, username, email, url_to_avatar FROM "USER" WHERE id=$1`

func (d UserDB) GetUserByID(ctx context.Context, ID int) (models.User, error) {
	var userInfo UserInfo

	err := d.Pool.QueryRow(ctx, getUserByIDQuery, ID).Scan(
		&userInfo.ID,
		&userInfo.Username,
		&userInfo.Email,
		&userInfo.ImageURL,
	)

	if err == pgx.ErrNoRows {
		return models.User{}, fmt.Errorf("%s: %w", models.LevelDB, models.ErrUserNotFound)
	}

	user := ToDomainUser(userInfo)
	return user, err
}
