//go:generate easyjson user.go
package models

//easyjson:json
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

type Subscription struct {
	SubscriberID int `json:"subscriber_id"`
	FollowsID    int `json:"subscribed_id"`
}
