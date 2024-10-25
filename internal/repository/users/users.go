package userRepository

import (
	"context"
	"errors"
	"time"

	"kudago/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserInfo struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	ImageURL  *string   `db:"url_to_avatar"`
	CreatedAt time.Time `db:"created_at"`
}

type UserDB struct {
	pool *pgxpool.Pool
}

func NewDB(pool *pgxpool.Pool) *UserDB {
	return &UserDB{
		pool: pool,
	}
}

func (d *UserDB) AddUser(ctx context.Context, user *models.User) (models.User, error) {
	rawQuery := `
		INSERT INTO "USER" (username, email, password_hash, url_to_avatar)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, url_to_avatar, created_at`

	var userInfo UserInfo
	err := d.pool.QueryRow(ctx, rawQuery,
		user.Username,
		user.Email,
		user.Password,
		user.ImageURL,
	).Scan(
		&userInfo.ID,
		&userInfo.Username,
		&userInfo.Email,
		&userInfo.ImageURL,
		&userInfo.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	newUser := toDomainUser(userInfo)
	return *&newUser, nil
}

func (d UserDB) CheckCredentials(ctx context.Context, username, password string) (models.User, error) {
	query := `
	SELECT id, username, email, created_at, url_to_avatar
	FROM "USER"
	WHERE username = $1 AND password_hash = $2`

	var userInfo UserInfo
	err := d.pool.QueryRow(ctx, query, username, password).Scan(
		&userInfo.ID,
		&userInfo.Username,
		&userInfo.Email,
		&userInfo.CreatedAt,
		&userInfo.ImageURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, err
	}
	user := toDomainUser(userInfo)
	return user, nil
}

func (d UserDB) GetUserByID(ctx context.Context, ID int) (models.User, error) {
	var userInfo UserInfo
	query := `SELECT id, username, email, url_to_avatar FROM "USER" WHERE id=$1`

	err := d.pool.QueryRow(ctx, query, ID).Scan(
		&userInfo.ID,
		&userInfo.Username,
		&userInfo.Email,
		&userInfo.ImageURL,
	)

	if err == pgx.ErrNoRows {
		return models.User{}, models.ErrUserNotFound
	}

	user := toDomainUser(userInfo)
	return user, err
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

func (d *UserDB) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT 1 FROM "USER" WHERE email=$1 LIMIT 1`

	var exists int
	err := d.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists > 0, nil
}

func (d *UserDB) UsernameExists(ctx context.Context, username string) (bool, error) {
	rawQuery := `SELECT 1 FROM "USER" WHERE username = $1 LIMIT 1`

	var exists int
	err := d.pool.QueryRow(ctx, rawQuery, username).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists > 0, nil
}
