package users

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Password    string `json:"-"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

	users["rvasily"] = User{
		ID:          0,
		Username:    "rvasily",
		Email:       "rvasily@example.com",
		DateOfBirth: "1990-05-15",
		Password:    "123",
	}

	users["ivanov"] = User{
		ID:          1,
		Username:    "ivanov",
		Email:       "ivanov@example.com",
		DateOfBirth: "1985-02-12",
		Password:    "234",
	}

	users["petrov"] = User{
		ID:          2,
		Username:    "petrov",
		Email:       "petrov@example.com",
		DateOfBirth: "1995-08-28",
		Password:    "345",
	}

	users["semenov"] = User{
		ID:          3,
		Username:    "semenov",
		Email:       "semenov@example.com",
		DateOfBirth: "1988-11-22",
		Password:    "456",
	}

	return users
}
