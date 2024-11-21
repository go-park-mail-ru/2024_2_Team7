package handlers

import (
	"kudago/internal/gateway"
	pb "kudago/internal/user/api"
)

type UserHandlers struct {
	Gateway *gateway.Gateway
}

type AuthResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ImageURL string `json:"image"`
}

type GetUsersResponse struct {
	Users []UserResponse `json:"users"`
}

func NewUserHandlers(gw *gateway.Gateway) *UserHandlers {
	return &UserHandlers{Gateway: gw}
}

func userToUserResponse(user *pb.User) UserResponse {
	return UserResponse{
		ID:       int(user.ID),
		Username: user.Username,
		Email:    user.Email,
		// ImageURL: user.ImageURL,
	}
}

func subscriptionToSubscriptionPB(subscription *pb.Subscription) pb.Subscription {
	return pb.Subscription{
		SubscriberID: int32(subscription.SubscriberID),
		FollowsID:    int32(subscription.SubscriberID),
	}
}

func writeUsersResponse(users []*pb.User, limit int) GetUsersResponse {
	resp := GetUsersResponse{Users: make([]UserResponse, 0, limit)}

	for _, user := range users {
		userResp := userToUserResponse(user)
		resp.Users = append(resp.Users, userResp)
	}
	return resp
}
