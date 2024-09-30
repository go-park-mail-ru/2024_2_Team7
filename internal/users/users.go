package users

func (d *UserDB) AddUser(user *User) {
	
	d.mu.Lock()
	user.ID = len(d.users)
	d.users[user.Username] = *user
	 d.mu.Unlock() 
}

func (d *UserDB) CheckCredentials(username, password string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	user, exists := d.users[username]
	if !exists || user.Password != password {
		return false
	}
	return true
}

func (d *UserDB) GetUser(username string) User {
	d.mu.RLock()
	user, _ := d.users[username]
	d.mu.RUnlock()
	return user
}

func (d *UserDB) UserExists(username string) bool {
	d.mu.RLock()
	
	_, exists := d.users[username]
	d.mu.RUnlock()
	return exists
}
