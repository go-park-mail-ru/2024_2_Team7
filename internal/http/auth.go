package http

import (
	"encoding/json"
	"net/http"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/models"
	"kudago/internal/repository"
)

const SessionToken = repository.SessionToken

type AuthHandler struct {
	Service AuthService
}

func NewAuthHandler(s AuthService) *AuthHandler {
	return &AuthHandler{
		Service: s,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(repository.SessionToken)
	if err == nil {
		_, authorized := h.Service.CheckSession(r.Context(), cookie.Value)
		if authorized {
			json.NewEncoder(w).Encode(httpErrors.ErrUserIsAuthorized)
			return
		}
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		json.NewEncoder(w).Encode(httpErrors.ErrInvalidRequest)
		return
	}

	if err := h.Service.Register(r.Context(), user); err != nil {
		json.NewEncoder(w).Encode(httpErrors.ErrUserAlreadyExists)
		return
	}

	user.Password = ""
	h.setSessionCookie(w, r, user.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(repository.SessionToken)
	if err == nil {
		_, authorized := h.Service.CheckSession(r.Context(), cookie.Value)
		if authorized {
			json.NewEncoder(w).Encode(httpErrors.ErrUserAlreadyLoggedIn)
			return
		}
	}

	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		json.NewEncoder(w).Encode(httpErrors.ErrInvalidRequest)
		return
	}

	if h.Service.CheckCredentials(r.Context(), creds) {
		user := h.Service.GetUser(r.Context(), creds.Username)
		h.setSessionCookie(w, r, creds.Username)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
		return
	}

	json.NewEncoder(w).Encode(httpErrors.ErrUnauthorized)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionToken)
	if err != nil {
		json.NewEncoder(w).Encode(httpErrors.ErrUnauthorized)
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
	w.WriteHeader(http.StatusOK)
	if err != nil {
		json.NewEncoder(w).Encode(httpErrors.ErrUnauthorized)
		return
	}
	session, authorized := h.Service.CheckSession(r.Context(), cookie.Value)

	if !authorized {
		json.NewEncoder(w).Encode(httpErrors.ErrUnauthorized)
		return
	}
	user := h.Service.GetUser(r.Context(), session.Username)
	user.Password = ""

	json.NewEncoder(w).Encode(user)
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
