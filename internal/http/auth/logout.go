package auth

import (
	"net/http"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"
)

// @Summary Выход из системы
// @Description Выход из аккаунта
// @Tags auth
// @Success 200
// @Failure 403 {object} httpErrors.HttpError "Forbidden"
// @Router /logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	err := h.service.DeleteSession(r.Context(), session.Token)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   models.SessionToken,
		MaxAge: -1, // Устанавливаем истекшее время, чтобы удалить cookie
	})

	w.WriteHeader(http.StatusOK)
}
