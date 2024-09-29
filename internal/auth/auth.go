package auth

import (
	"encoding/json"
	"net/http"

	"kudago/internal/users"
	"kudago/session"
)

type AuthHandler struct {
	UserDB    users.UserDB
	SessionDb session.SessionDB
}

const SessionToken = session.SessionToken

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	session, authorized := h.SessionDb.CheckSession(r)
	if authorized {
		http.Error(w, "Alredy logged in", http.StatusForbidden)
		return
	}

	var user users.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if h.UserDB.UserExists(user.Username) {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	h.UserDB.AddUser(user)
	session = h.SessionDb.CreateSession(user.Username)
	http.SetCookie(w, &http.Cookie{
		Name:     SessionToken,
		Value:    session.Token,
		Expires:  session.Expires,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusCreated)
	creds := h.UserDB.GetCredentials(user.Username)
	json.NewEncoder(w).Encode(creds)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	_, authorized := h.SessionDb.CheckSession(r)
	if authorized {
		http.Error(w, "Alredy logged in", http.StatusForbidden)
		return
	}

	var creds users.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if h.UserDB.CheckCredentials(creds.Username, creds.Password) {
		user := h.UserDB.GetCredentials(creds.Username)
		session := h.SessionDb.CreateSession(creds.Username)
		http.SetCookie(w, &http.Cookie{
			Name:     SessionToken,
			Value:    session.Token,
			Expires:  session.Expires,
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.SessionDb.DeleteSession(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:   SessionToken,
		MaxAge: -1, // Устанавливаем истекшее время, чтобы удалить cookie
	})

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	session, authorized := h.SessionDb.CheckSession(r)
	if !authorized {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}
