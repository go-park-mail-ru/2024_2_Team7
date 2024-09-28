package users

func (d *UserDB) AddUser(user User) {
	d.users[user.Username] = user
}

func (d *UserDB) CheckCredentials(username, password string) bool {
	user, exists := d.users[username]
	if !exists || user.Password != password {
		return false
	}
	return true
}

func (d *UserDB) GetUser(username string) User {
	user, _ := d.users[username]
	return user
}
