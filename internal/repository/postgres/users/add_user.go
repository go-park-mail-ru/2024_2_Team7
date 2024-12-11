package userRepository

import (
	"context"
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/argon2"

	"kudago/internal/models"
)

const addUserQuery = `
	INSERT INTO "USER" (username, email, password_hash, url_to_avatar)
	VALUES ($1, $2, $3, $4)
	RETURNING id,  created_at`

func (d *UserDB) AddUser(ctx context.Context, user models.User) (models.User, error) {
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return models.User{}, fmt.Errorf("failed to generate salt: %w", err)
	}
	hashedPass := argon2.IDKey([]byte(user.Password), []byte(salt), 1, 64*1024, 4, 32)
	passwordHash := append(salt, hashedPass...)

	var userInfo UserInfo
	err := d.pool.QueryRow(ctx, addUserQuery,
		user.Username,
		user.Email,
		passwordHash,
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
	newUser := toDomainUser(userInfo)
	return newUser, nil
}
