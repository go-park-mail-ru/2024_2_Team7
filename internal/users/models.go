package users

import "sync"

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDB struct {
	users map[string]User
	mu    *sync.RWMutex
}

func NewUserDB() *UserDB {
	users := createUserMapWithDefaultValues()

	return &UserDB{
		users: users,
		mu:       &sync.RWMutex{}, 
	}
}

func createUserMapWithDefaultValues() map[string]User {
	users := make(map[string]User)

	users["rvasily"] = User{
		ID:          0,
		Username:    "rvasily",
		Email:       "rvasily@example.com",
		Password:    "123",
	}

	users["ivanov"] = User{
		ID:          1,
		Username:    "ivanov",
		Email:       "ivanov@example.com",
		Password:    "234",
	}

	users["petrov"] = User{
		ID:          2,
		Username:    "petrov",
		Email:       "petrov@example.com",
		Password:    "345",
	}

	users["semenov"] = User{
		ID:          3,
		Username:    "semenov",
		Email:       "semenov@example.com",
		Password:    "456",
	}

	return users
}
