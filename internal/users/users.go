package users

import "errors"

var users = make(map[string]User)

func AddUser(user User) {
	users[user.Username] = user
}

func CheckCredentials(username, password string) bool {
	user, exists := users[username]
	if !exists || user.Password != password {
		return false
	}
	return true
}

func GetUser(username string) (User, error) {
	user, exists := users[username]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}
