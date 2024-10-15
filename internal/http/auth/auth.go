package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"kudago/internal/models"
	"kudago/internal/repository"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
)

const SessionToken = repository.SessionToken

type AuthHandler struct {
	Service AuthService
	decoder *schema.Decoder
}

type AuthService interface {
	CheckSession(ctx context.Context, cookie string) (*models.Session, bool)
	GetUser(ctx context.Context, username string) models.User
	CheckCredentials(ctx context.Context, creds models.Credentials) bool
	Register(ctx context.Context, user models.User) (models.User, error)
	CreateSession(ctx context.Context, username string) *models.Session
	DeleteSession(ctx context.Context, username string)
}

type RegisterRequest struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required,alphanum,length(3|50)"`
}

type CredentialsBucket struct {
	Username string `json:"username" valid:"required,alphanum,length(3|50)"`
	Password string `json:"password" valid:"required,alphanum,length(3|50)"`
}

type AuthRequest CredentialsBucket

type AuthResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewAuthHandler(s AuthService) *AuthHandler {
	return &AuthHandler{
		Service: s,
		decoder: schema.NewDecoder(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(repository.SessionToken)
	if err == nil {
		_, authorized := h.Service.CheckSession(r.Context(), cookie.Value)
		if authorized {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(errUserIsAuthorized)
			return
		}
	}

	var req RegisterRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errInvalidRequest)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errInvalidFields)
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err = h.Service.Register(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errInvalidData)
		return
	}

	h.setSessionCookie(w, r, user.Username)

	userResponse := userToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(repository.SessionToken)
	if err == nil {
		_, authorized := h.Service.CheckSession(r.Context(), cookie.Value)
		if authorized {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(errUserAlreadyLoggedIn)
			return
		}
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errInvalidRequest)
		return
	}

	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errInvalidFields)
		return
	}

	creds := models.Credentials{
		Username: req.Username,
		Password: req.Password,
	}

	if h.Service.CheckCredentials(r.Context(), creds) {
		user := h.Service.GetUser(r.Context(), creds.Username)
		h.setSessionCookie(w, r, creds.Username)
		w.WriteHeader(http.StatusOK)
		userResponse := userToUserResponse(user)

		resp := AuthResponse{
			User: userResponse,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(errUnauthorized)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionToken)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(errUnauthorized)
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
	cookie, err := r.Cookie(SessionToken)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(errUnauthorized)
		return
	}

	session, authorized := h.Service.CheckSession(r.Context(), cookie.Value)

	if !authorized {
		json.NewEncoder(w).Encode(errUnauthorized)
		return
	}
	user := h.Service.GetUser(r.Context(), session.Username)
	userResponse := userToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, r *http.Request, username string) {
	session := h.Service.CreateSession(r.Context(), username)
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
