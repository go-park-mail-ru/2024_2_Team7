package gateway

import (
	"context"
	"net/http"
	"time"

	pb "kudago/internal/auth/api"

	"kudago/internal/models"
)

type AuthResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ImageURL string `json:"image"`
}

func userToUserResponse(user *pb.User) AuthResponse {
	resp := AuthResponse{
		User: UserResponse{
			ID:       int(user.ID),
			Username: user.Username,
			Email:    user.Email,
			ImageURL: user.AvatarUrl,
		},
	}
	return resp
}

type sessionKeyType struct{}

var sessionKey sessionKeyType

func SetSessionInContext(ctx context.Context, session models.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

func (g *Gateway) setSessionCookie(w http.ResponseWriter, r *http.Request, ID int) error {
	req := &pb.CreateSessionRequest{ID: int32(ID)}

	session, err := g.authClient.CreateSession(r.Context(), req)
	if err != nil {
		return models.ErrInternal
	}

	expires, err := time.Parse(time.RFC3339, session.Expires)
	if err != nil {
		return models.ErrInternal
	}

	http.SetCookie(w, &http.Cookie{
		Name:     models.SessionToken,
		Value:    session.Token,
		Expires:  expires,
		HttpOnly: true,
	})
	return nil
}
