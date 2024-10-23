package userRepository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"kudago/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	userColumns = []string{"id", "username", "email", "password_hash", "created_at", "url_to_avatar"}
	userTable   = `"USER"`
)

type UserInfo struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
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
	insertColumns := []string{"username", "email", "password_hash", "url_to_avatar"}

	query := sq.Insert(userTable).
		Columns(insertColumns...).
		Values(user.Username, user.Email, user.Password, user.ImageURL).
		Suffix("RETURNING id, username, email, password_hash, url_to_avatar, created_at").
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return models.User{}, err
	}

	var userInfo UserInfo
	err = d.pool.QueryRow(ctx, rawQuery, args...).Scan(
		&userInfo.ID,
		&userInfo.Username,
		&userInfo.Email,
		&userInfo.Password,
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
	query := sq.Select(userColumns...).
		From(userTable).
		Where(sq.Eq{"username": username, "password_hash": password}).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return models.User{}, err
	}

	var userInfo UserInfo
	err = d.pool.QueryRow(ctx, rawQuery, args...).Scan(
		&userInfo.ID,
		&userInfo.Username,
		&userInfo.Email,
		&userInfo.Password,
		&userInfo.CreatedAt,
		&userInfo.ImageURL,
	)
	fmt.Println(err, errors.Is(err, pgx.ErrNoRows))
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
	query := sq.Select(userColumns...).
		From(userTable).
		Where(sq.Eq{"id": ID}).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return models.User{}, err
	}

	err = d.pool.QueryRow(ctx, rawQuery, args...).Scan(
		&userInfo.ID,
		&userInfo.Username,
		&userInfo.Email,
		&userInfo.Password,
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
		Password: user.Password,
		ImageURL: imageURL,
	}
}

func (d *UserDB) EmailExists(ctx context.Context, email string) (bool, error) {
	query := sq.Select("1").
		From(userTable).
		Where(sq.Eq{"email": email}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = d.pool.QueryRow(ctx, rawQuery, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists > 0, nil
}

func (d *UserDB) UsernameExists(ctx context.Context, username string) (bool, error) {
	query := sq.Select("1").
		From(userTable).
		Where(sq.Eq{"username": username}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = d.pool.QueryRow(ctx, rawQuery, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists > 0, nil
}
