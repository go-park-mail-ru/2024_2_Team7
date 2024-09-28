package users

type User struct {
	Username string
	Password string
}

type UserDB struct {
	users map[string]User
}

func NewUserDB() *UserDB {
	users := createUserMapWithDefaultValues()

	return &UserDB{
		users: users,
	}
}

func createUserMapWithDefaultValues() map[string]User {
	users := make(map[string]User)
	users["rvasily"] = User{Password: "123"}
	users["ivanov"] = User{Password: "123"}
	users["petrov"] = User{Password: "123"}
	users["semenov"] = User{Password: "123"}
	return users
}
