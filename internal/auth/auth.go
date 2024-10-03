package auth

import (
	"encoding/json"
	"net/http"

	"kudago/internal/users"
	"kudago/session"
)

type Handler struct {
	UserDB    users.UserDB
	SessionDb session.SessionDB
}

type authError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

const SessionToken = session.SessionToken

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	_, authorized := h.SessionDb.CheckSession(r)
	if authorized {
		err := &authError{
			Message: "User is authorized",
			Code:    http.StatusForbidden,
		}
		json.NewEncoder(w).Encode(err)
		return
	}

	var user users.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		err := &authError{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		}
		json.NewEncoder(w).Encode(err)
		return
	}

	if h.UserDB.UserExists(user.Username) {
		err := &authError{
			Message: "User alresdy exists",
			Code:    http.StatusConflict,
		}
		json.NewEncoder(w).Encode(err)
		return
	}

	if err := h.UserDB.AddUser(&user); err != nil {
		err := &authError{
			Message: "Email is already used",
			Code:    http.StatusConflict,
		}
		json.NewEncoder(w).Encode(err)
		return
	}
	user.Password = ""
	h.setSessionCookie(w, user.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	_, authorized := h.SessionDb.CheckSession(r)
	if authorized {
		err := &authError{
			Message: "Already logged in",
			Code:    http.StatusForbidden,
		}
		json.NewEncoder(w).Encode(err)
		return
	}

	var creds users.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		err := &authError{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		}
		json.NewEncoder(w).Encode(err)
		return
	}

	if h.UserDB.CheckCredentials(creds.Username, creds.Password) {
		user := h.UserDB.GetUser(creds.Username)
		user.Password = ""
		h.setSessionCookie(w, creds.Username)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
		return
	}
	err := &authError{
		Message: "Unauthorized",
		Code:    http.StatusUnauthorized,
	}
	json.NewEncoder(w).Encode(err)

}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionToken)
	if err != nil {
		err := &authError{
			Message: "Unauthorized",
			Code:    http.StatusUnauthorized,
		}
		json.NewEncoder(w).Encode(err)
		return
	}

	h.SessionDb.DeleteSession(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:   SessionToken,
		MaxAge: -1, // Устанавливаем истекшее время, чтобы удалить cookie
	})

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CheckSession(w http.ResponseWriter, r *http.Request) {
	session, authorized := h.SessionDb.CheckSession(r)
	w.WriteHeader(http.StatusOK)

	if !authorized {
		err := &authError{
			Message: "User is not authorized",
			Code:    http.StatusUnauthorized,
		}
		json.NewEncoder(w).Encode(err)
		return
	}
	user := h.UserDB.GetUser(session.Username)
	user.Password = ""

	json.NewEncoder(w).Encode(user)
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, username string) {
	session := h.SessionDb.CreateSession(username)
	http.SetCookie(w, &http.Cookie{
		Name:     SessionToken,
		Value:    session.Token,
		Expires:  session.Expires,
		HttpOnly: true,
	})
}
