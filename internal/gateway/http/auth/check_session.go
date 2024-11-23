package auth

import (
	"net/http"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
)

// @Summary Проверка сессии
// @Description Возвращает информацию о пользователе, если сессия активна
// @Tags auth
// @Success 200 {object} AuthResponse
// @Router /session [get]
func (h *AuthHandler) CheckSession(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusOK, httpErrors.ErrUnauthorized)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
		return
	}

	userResponse := userToUserResponse(user)

	resp := AuthResponse{
		User: userResponse,
	}

	utils.WriteResponse(w, http.StatusOK, resp)
}
