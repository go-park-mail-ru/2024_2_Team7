package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"kudago/internal/models"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
)

const (
	SessionToken = "session_token"
	SessionKey   = "session"
)

type AuthHandler struct {
	Service AuthService
	decoder *schema.Decoder
}

type AuthService interface {
	CheckSession(ctx context.Context, cookie string) (*models.Session, bool)
	GetUserByUsername(ctx context.Context, username string) models.User
	GetUserByID(ctx context.Context, ID int) models.User
	CheckCredentials(ctx context.Context, creds models.Credentials) bool
	Register(ctx context.Context, user models.User) (models.User, error)
	CreateSession(ctx context.Context, ID int) *models.Session
	DeleteSession(ctx context.Context, token string)
}

type RegisterRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
}

type AuthRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Password string `json:"password" valid:"password,required,length(3|50)"`
}

type AuthResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ProfileResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ImageURL string `json:"image"`
}

type ValidationErrResponse struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func init() {
	govalidator.TagMap["password"] = govalidator.Validator(func(str string) bool {
		regex := `^[a-zA-Z0-9+\-*/.;=\]\[\}\{\?]+$`
		match, _ := regexp.MatchString(regex, str)
		return match
	})
}

func NewAuthHandler(s AuthService) *AuthHandler {
	return &AuthHandler{
		Service: s,
		decoder: schema.NewDecoder(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	_, ok := getSessionFromContext(r.Context())
	if ok {
		writeResponse(w, http.StatusForbidden, errUserIsAuthorized)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, errInvalidFields)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		processValidationErrors(w, err)
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err = h.Service.Register(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUsernameIsUsed):
			writeResponse(w, http.StatusConflict, errUsernameIsAlredyTaken)
		case errors.Is(err, models.ErrEmailIsUsed):
			writeResponse(w, http.StatusConflict, errEmailIsAlredyTaken)
		default:
			writeResponse(w, http.StatusInternalServerError, errInternal)
		}
		return
	}

	h.setSessionCookie(w, r, user.ID)

	userResponse := userToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	writeResponse(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	_, ok := getSessionFromContext(r.Context())
	if ok {
		writeResponse(w, http.StatusForbidden, errUserAlreadyLoggedIn)
		return
	}
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, errInvalidFields)
		return
	}

	_, err := govalidator.ValidateStruct(&req)
	if err != nil {
		processValidationErrors(w, err)
		return
	}

	creds := models.Credentials{
		Username: req.Username,
		Password: req.Password,
	}

	if h.Service.CheckCredentials(r.Context(), creds) {
		user := h.Service.GetUserByUsername(r.Context(), creds.Username)
		h.setSessionCookie(w, r, user.ID)
		userResponse := userToUserResponse(user)

		resp := AuthResponse{
			User: userResponse,
		}

		writeResponse(w, http.StatusOK, resp)
		return
	}
	writeResponse(w, http.StatusForbidden, errUnauthorized)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionToken)
	if err != nil {
		writeResponse(w, http.StatusForbidden, errUnauthorized)
		return
	}

	h.Service.DeleteSession(r.Context(), cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:   SessionToken,
		MaxAge: -1, // Устанавливаем истекшее время, чтобы удалить cookie
	})

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) CheckSession(w http.ResponseWriter, r *http.Request) {
	session, ok := getSessionFromContext(r.Context())
	if !ok {
		writeResponse(w, http.StatusOK, errUnauthorized)
		return
	}

	user := h.Service.GetUserByID(r.Context(), session.UserID)
	userResponse := userToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	writeResponse(w, http.StatusOK, resp)
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	session, ok := getSessionFromContext(r.Context())
	if !ok {
		writeResponse(w, http.StatusOK, errUnauthorized)
		return
	}

	user := h.Service.GetUserByID(r.Context(), session.UserID)
	userResponse := userToProfileResponse(user)

	writeResponse(w, http.StatusOK, userResponse)
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, r *http.Request, ID int) {
	session := h.Service.CreateSession(r.Context(), ID)
	http.SetCookie(w, &http.Cookie{
		Name:     SessionToken,
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
	}
}

func userToProfileResponse(user models.User) ProfileResponse {
	return ProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ImageURL: user.ImageURL,
	}
}

func getSessionFromContext(ctx context.Context) (*models.Session, bool) {
	session, ok := ctx.Value(SessionKey).(*models.Session)
	return session, ok
}

func writeResponse(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func processValidationErrors(w http.ResponseWriter, err error) {
	errors := strings.Split(err.Error(), ";")
	resp := ValidationErrResponse{}

	for _, err := range errors {
		colonIndex := strings.Index(err, ":")

		if colonIndex == -1 {
			continue
		}

		field := err[:colonIndex]
		errorMsg := err[colonIndex+2:]

		valErr := ValidationError{
			Field: field,
			Error: errorMsg,
		}
		resp.Errors = append(resp.Errors, valErr)
	}
	writeResponse(w, http.StatusUnauthorized, resp)
}
