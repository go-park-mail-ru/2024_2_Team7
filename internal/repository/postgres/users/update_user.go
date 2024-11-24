package userRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
)

const updateUserQuery = `
	UPDATE "USER"
	SET 
		username = COALESCE($2, username), 
		email = COALESCE($3, email), 
		URL_to_avatar = COALESCE($4, URL_to_avatar), 
		modified_at = NOW()
	WHERE id = $1 
	RETURNING id, username, email, URL_to_avatar
`

func (db *UserDB) UpdateUser(ctx context.Context, updatedUser models.User) (models.User, error) {
	var existingID int
	row := db.pool.QueryRow(ctx, getUserByEmailOrUsernameQuery, updatedUser.Email, updatedUser.Username)
	if err := row.Scan(&existingID); err != pgx.ErrNoRows && err != nil {
		return models.User{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	if existingID != 0 {
		return models.User{}, models.ErrEmailIsUsed
	}

	var user models.User
	err := db.pool.QueryRow(ctx, updateUserQuery,
		updatedUser.ID,
		nilIfEmpty(updatedUser.Username),
		nilIfEmpty(updatedUser.Email),
		nilIfEmpty(updatedUser.ImageURL),
	).Scan(&user.ID, &user.Username, &user.Email, &user.ImageURL)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	return user, nil
}
