package userRepository

import (
	"context"
	"time"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type UserInfo struct {
	ID         int       `db:"id"`
	Username   string    `db:"username"`
	Email      string    `db:"email"`
	ImageURL   *string   `db:"url_to_avatar"`
	CreatedAt  time.Time `db:"created_at"`
	ModifiedAt time.Time `db:"modified_at"`
}

type UserDB struct {
	pool *pgxpool.Pool
}

func NewDB(pool *pgxpool.Pool) *UserDB {
	return &UserDB{
		pool: pool,
	}
}

const checkUsernameDuplicateQuery = `
	SELECT id 
	FROM "USER"
	WHERE username = $1  AND id != $2
`

func (db *UserDB) CheckUsername(ctx context.Context, username string, ID int) (bool, error) {
	var exists int
	err := db.pool.QueryRow(ctx, checkUsernameDuplicateQuery, username, ID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, models.LevelDB)
	}

	return true, nil
}

const checkEmailDuplicateQuery = `
	SELECT id 
	FROM "USER"
	WHERE email = $1  AND id != $2
`

func (db *UserDB) CheckEmail(ctx context.Context, email string, ID int) (bool, error) {
	var exists int
	err := db.pool.QueryRow(ctx, checkEmailDuplicateQuery, email, ID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, models.LevelDB)
	}

	return true, nil
}

func nilIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func toDomainUser(user UserInfo) models.User {
	var imageURL string
	if user.ImageURL == nil {
		imageURL = ""
	} else {
		imageURL = *user.ImageURL
	}

	return models.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ImageURL: imageURL,
	}
}
