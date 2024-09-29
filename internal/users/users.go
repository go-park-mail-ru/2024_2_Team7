package users

import "fmt"

func (d *UserDB) AddUser(user User) {
	user.ID = fmt.Sprintf("%d", len(d.users))
	d.users[user.Username] = user
}

func (d *UserDB) CheckCredentials(username, password string) bool {
	user, exists := d.users[username]
	if !exists || user.Password != password {
		return false
	}
	return true
}

func (d *UserDB) GetCredentials(username string) Credentials {
	user := d.users[username]
	return Credentials{
		ID:          user.ID,
		Username:    user.Username,
		DateOfBirth: user.DateOfBirth,
		Email:       user.Email,
	}
}

func (d *UserDB) GetUser(username string) User {
	user, _ := d.users[username]
	return user
}

func (d *UserDB) UserExists(username string) bool {
	_, exists := d.users[username]
	return exists
}
