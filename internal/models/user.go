package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ImageURL string `json:"image"`
}

type NewUserData struct {
	User  User
	Media MediaFile
}
