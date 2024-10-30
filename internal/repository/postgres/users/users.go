package userRepository

import (
	"context"
	"errors"
	"time"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	addUserQuery = `
	INSERT INTO "USER" (username, email, password_hash, url_to_avatar)
	VALUES ($1, $2, $3, $4)
	RETURNING id,  created_at`

	checkCredentialsQuery = `
	SELECT id, username, email, created_at, url_to_avatar
	FROM "USER"
	WHERE username = $1 AND password_hash = $2`

	getUserByIDQuery = `SELECT id, username, email, url_to_avatar FROM "USER" WHERE id=$1`

	getUserByEmailOrUsernameQuery = `SELECT 1 FROM "USER" WHERE email=$1 OR username = $2 LIMIT 1`
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

func (d *UserDB) AddUser(ctx context.Context, user models.User) (models.User, error) {
	var userInfo UserInfo
	err := d.pool.QueryRow(ctx, addUserQuery,
		user.Username,
		user.Email,
		user.Password,
		user.ImageURL,
	).Scan(
		&userInfo.ID,
		&userInfo.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	userInfo.Username = user.Username
	userInfo.Email = user.Email
	userInfo.ImageURL = &user.ImageURL
	newUser := toDomainUser(userInfo)
	return newUser, nil
}

func (d UserDB) CheckCredentials(ctx context.Context, username, password string) (models.User, error) {
	var userInfo UserInfo
	err := d.pool.QueryRow(ctx, checkCredentialsQuery, username, password).Scan(
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

	err := d.pool.QueryRow(ctx, getUserByIDQuery, ID).Scan(
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

func (d *UserDB) UserExists(ctx context.Context, username, email string) (bool, error) {
	var exists int
	err := d.pool.QueryRow(ctx, getUserByEmailOrUsernameQuery, email, username).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
