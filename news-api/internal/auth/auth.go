package auth

import (
	"encoding/json"
	"net/http"
	"news-api/internal/users"
	"news-api/pkg/session"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user users.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	users.AddUser(user)
	w.WriteHeader(http.StatusCreated)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds users.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if users.CheckCredentials(creds.Username, creds.Password) {
		sessionToken := session.CreateSession(creds.Username)
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: session.GetExpiration(),
		})
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем cookie сессии
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Удаляем сессию из хранилища сессий
	session.DeleteSession(cookie.Value)

	// Удаляем cookie у пользователя
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		MaxAge: -1, // Устанавливаем истекшее время, чтобы удалить cookie
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
