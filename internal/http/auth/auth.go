//go:generate mockgen -source ./auth.go -destination=./mocks/auth.go -package=mocks

package auth

import (
	"context"
	"net/http"
	"regexp"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
)

type AuthHandler struct {
	service AuthService
	logger  *logger.Logger
}

type AuthService interface {
	CheckSession(ctx context.Context, cookie string) (models.Session, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	UpdateUser(ctx context.Context, data models.NewUserData) (models.User, error)
	CheckCredentials(ctx context.Context, creds models.Credentials) (models.User, error)
	Register(ctx context.Context, registerDTO models.NewUserData) (models.User, error)
	CreateSession(ctx context.Context, ID int) (models.Session, error)
	DeleteSession(ctx context.Context, token string) error
	Subscribe(ctx context.Context, subscription models.Subscription) error
	Unsubscribe(ctx context.Context, subscription models.Subscription) error
	GetSubscriptions(ctx context.Context, ID int) ([]models.User, error)
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

var validPasswordRegex = regexp.MustCompile(`^[a-zA-Z0-9+\-*/.;=\]\[\}\{\?]+$`)

func init() {
	govalidator.TagMap["password"] = govalidator.Validator(func(str string) bool {
		return validPasswordRegex.MatchString(str)
	})
}

func NewAuthHandler(s AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		service: s,
		logger:  logger,
	}
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, r *http.Request, ID int) {
	session, err := h.service.CreateSession(r.Context(), ID)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     models.SessionToken,
		Value:    session.Token,
		Expires:  session.Expires,
		HttpOnly: true,
	})
}

func userToUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ImageURL: user.ImageURL,
	}
}

func writeUsersResponse(users []models.User, limit int) GetUsersResponse {
	resp := GetUsersResponse{make([]UserResponse, 0, limit)}

	for _, user := range users {
		userResp := userToUserResponse(user)
		resp.Users = append(resp.Users, userResp)
	}
	return resp
}
