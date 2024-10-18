package userRepository

import (
	"context"
	"errors"
	"strings"
	"sync"

	"kudago/internal/models"
)

var errEmailIsUsed = errors.New("Email is already used")

type UserDB struct {
	users map[string]models.User
	mu    *sync.RWMutex
}

func NewDB() *UserDB {
	users := createUserMapWithDefaultValues()

	return &UserDB{
		users: users,
		mu:    &sync.RWMutex{},
	}
}

func (d *UserDB) AddUser(ctx context.Context, user *models.User) (models.User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, u := range d.users {
		if strings.ToLower(user.Email) == strings.ToLower(u.Email) {
			return models.User{}, errEmailIsUsed
		}
	}
	user.ID = len(d.users)
	d.users[user.Username] = *user
	return *user, nil
}

func (d UserDB) CheckCredentials(ctx context.Context, username, password string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	user, exists := d.users[username]
	if !exists || user.Password != password {
		return false
	}
	return true
}

func (d UserDB) GetUserByUsername(ctx context.Context, username string) models.User {
	d.mu.RLock()
	user, _ := d.users[username]
	d.mu.RUnlock()
	return user
}

func (d UserDB) GetUserByID(ctx context.Context, ID int) models.User {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, user := range d.users {
		if user.ID == ID {
			return user
		}
	}
	return models.User{}
}

func (d UserDB) UserExists(ctx context.Context, username string) bool {
	d.mu.RLock()
	_, exists := d.users[username]
	d.mu.RUnlock()
	return exists
}

func createUserMapWithDefaultValues() map[string]models.User {
	users := make(map[string]models.User)

	users["rvasily"] = models.User{
		ID:       0,
		Username: "rvasily",
		Email:    "rvasily@example.com",
		ImageURL: "/static/images/profile1.jpeg",
		Password: "123",
	}

	users["ivanov"] = models.User{
		ID:       1,
		Username: "ivanov",
		Email:    "ivanov@example.com",
		Password: "234",
	}

	users["petrov"] = models.User{
		ID:       2,
		Username: "petrov",
		Email:    "petrov@example.com",
		Password: "345",
	}

	users["semenov"] = models.User{
		ID:       3,
		Username: "semenov",
		Email:    "semenov@example.com",
		Password: "456",
	}

	return users
}
