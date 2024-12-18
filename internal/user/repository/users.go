package userRepository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"

	"kudago/internal/models"
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
	Pool Pool
}

type Pool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func NewDB(pool Pool) *UserDB {
	return &UserDB{
		Pool: pool,
	}
}

func NilIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func ToDomainUser(user UserInfo) models.User {
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
