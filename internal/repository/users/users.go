package userRepository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"kudago/internal/db"
	"kudago/internal/models"
	"log"
)

var errEmailIsUsed = errors.New("Email is already used")

type UserDB struct {
	pool *pgxpool.Pool
}

func NewDB() *UserDB {
	return &UserDB{
		pool: db.GetDB(),
	}
}

func (d *UserDB) AddUser(ctx context.Context, user *models.User) (models.User, error) {
	query := `
        INSERT INTO "USER" (name, email, password, created_at)
        VALUES ($1, $2, $3, CURRENT_DATE)
        RETURNING id
    `
	err := d.pool.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.ImageURL).Scan(&user.ID)
	if err != nil {
		if err.Error() == "unique_violation" {
			return models.User{}, errEmailIsUsed
		}
		return models.User{}, err
	}
	return *user, nil
}

func (d UserDB) CheckCredentials(ctx context.Context, username, password string) bool {
	query := `
        SELECT COUNT(*) FROM "USER" WHERE name=$1 AND password=$2
    `
	var count int
	err := d.pool.QueryRow(ctx, query, username, password).Scan(&count)
	return err == nil && count > 0
}

func (d UserDB) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	query := `
        SELECT id, name, email, password, URL_to_avatar FROM "USER" WHERE name=$1
    `
	err := d.pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ImageURL,
	)
	if err == pgx.ErrNoRows {
		return models.User{}, errors.New("user not found")
	}
	return user, err
}

func (d UserDB) GetUserByID(ctx context.Context, ID int) (models.User, error) {
	var user models.User
	query := `
        SELECT id, name, email, password, URL_to_avatar FROM "USER" WHERE id=$1
    `
	err := d.pool.QueryRow(ctx, query, ID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ImageURL,
	)
	if err == pgx.ErrNoRows {
		return models.User{}, errors.New("user not found")
	}
	return user, err
}

func (d UserDB) UserExists(ctx context.Context, username string) bool {
	query := `
        SELECT EXISTS(SELECT 1 FROM "USER" WHERE name=$1)
    `
	var exists bool
	err := d.pool.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		log.Fatal("error in service/auth/userExists")
	}
	return exists
}
