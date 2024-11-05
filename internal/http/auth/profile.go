package auth

import (
	"net/http"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"
)

type ProfileResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ImageURL string `json:"image"`
}

// @Summary Профиль пользователя
// @Description Возвращает информацию о профиле текущего пользователя
// @Tags profile
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Router /profile [get]
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
		return
	}
	user, err := h.service.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
		return
	}

	userResponse := userToProfileResponse(user)

	utils.WriteResponse(w, http.StatusOK, userResponse)
}

func userToProfileResponse(user models.User) ProfileResponse {
	return ProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		ImageURL: user.ImageURL,
	}
}
