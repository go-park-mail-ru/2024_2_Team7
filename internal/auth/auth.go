package auth

import (
	"encoding/json"
	"net/http"

	"kudago/internal/events"
	"kudago/internal/users"
	"kudago/session"
)

type Handler struct {
	UserDB    users.UserDB
	SessionDb session.SessionDB
	EventDB   events.EventDB
}

const sessionToken = "session_token"

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user users.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	h.UserDB.AddUser(user)
	session := h.SessionDb.CreateSession(user.Username)
	http.SetCookie(w, &http.Cookie{
		Name:    sessionToken,
		Value:   session.Token,
		Expires: session.Expires,
	})
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds users.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if h.UserDB.CheckCredentials(creds.Username, creds.Password) {
		user := h.UserDB.GetUser(creds.Username)
		session := h.SessionDb.CreateSession(creds.Username)
		http.SetCookie(w, &http.Cookie{
			Name:    sessionToken,
			Value:   session.Token,
			Expires: session.Expires,
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем cookie сессии
	cookie, err := r.Cookie(sessionToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.SessionDb.DeleteSession(cookie.Value)

	// Удаляем cookie у пользователя
	http.SetCookie(w, &http.Cookie{
		Name:   sessionToken,
		MaxAge: -1, // Устанавливаем истекшее время, чтобы удалить cookie
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("logged out")
}
