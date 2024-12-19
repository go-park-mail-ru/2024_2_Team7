package userRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const createUserQuery = `
	INSERT INTO "USER" (username, email, password_hash, url_to_avatar)
	VALUES ($1, $2, $3, $4)
	RETURNING id,  created_at`

func (d *UserDB) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var userInfo UserInfo
	err := d.Pool.QueryRow(ctx, createUserQuery,
		user.Username,
		user.Email,
		user.Password,
		user.ImageURL,
	).Scan(
		&userInfo.ID,
		&userInfo.CreatedAt,
	)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	userInfo.Username = user.Username
	userInfo.Email = user.Email
	userInfo.ImageURL = &user.ImageURL

	newUser := ToDomainUser(userInfo)
	return newUser, nil
}
