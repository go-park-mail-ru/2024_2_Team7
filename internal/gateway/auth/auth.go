package handlers

import (
	"net/http"
	"regexp"
	"time"

	pb "kudago/internal/auth/api"
	"kudago/internal/gateway"

	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

type AuthHandlers struct {
	Gateway *gateway.Gateway
}

var validPasswordRegex = regexp.MustCompile(`^[a-zA-Z0-9+\-*/.;=\]\[\}\{\?]+$`)

func init() {
	govalidator.TagMap["password"] = govalidator.Validator(func(str string) bool {
		return validPasswordRegex.MatchString(str)
	})
}

func NewAuthHandlers(gw *gateway.Gateway) *AuthHandlers {
	return &AuthHandlers{Gateway: gw}
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

func (h *AuthHandlers) setSessionCookie(w http.ResponseWriter, r *http.Request, ID int) error {
	req := &pb.CreateSessionRequest{ID: int32(ID)}

	session, err := h.Gateway.AuthService.CreateSession(r.Context(), req)
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
