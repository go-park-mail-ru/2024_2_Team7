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

const SessionToken = session.SessionToken

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	_, authorized := h.SessionDb.CheckSession(r)
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
	h.setSessionCookie(w, user.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	_, authorized := h.SessionDb.CheckSession(r)
	if authorized {
		http.Error(w, "Alredy logged in", http.StatusForbidden)
		return
	}

	var creds users.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if h.UserDB.CheckCredentials(creds.Username, creds.Password) {
		user := h.UserDB.GetUser(creds.Username)
		h.setSessionCookie(w, creds.Username)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) CheckSession(w http.ResponseWriter, r *http.Request) {
	session, authorized := h.SessionDb.CheckSession(r)
	if !authorized {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	user:=h.UserDB.GetUser(session.Username)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, username string){
	session := h.SessionDb.CreateSession(username)
	http.SetCookie(w, &http.Cookie{
		Name:     SessionToken,
		Value:    session.Token,
		Expires:  session.Expires,
		HttpOnly: true,
	})
}