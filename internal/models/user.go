package models

import "mime/multipart"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ImageURL string `json:"image"`
}

type RegisterDTO struct {
	User   User
	Header multipart.FileHeader
	File   multipart.File
}
