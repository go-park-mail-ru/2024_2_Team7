package handlers

import (
	"net/http"
	"strconv"

	pb "kudago/internal/user/api"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"

	"github.com/gorilla/mux"
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
// @Router /profile/{id} [get]
func (h *UserHandlers) Profile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidID)
		return
	}

	user, err := h.UserService.GetUserByID(r.Context(), &pb.GetUserByIDRequest{ID: int32(id)})
	if err != nil {
		utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrUserNotFound)
		return
	}

	userResponse := userToProfileResponse(user)

	utils.WriteResponse(w, http.StatusOK, userResponse)
}

func userToProfileResponse(user *pb.User) ProfileResponse {
	return ProfileResponse{
		ID:       int(user.ID),
		Username: user.Username,
		Email:    user.Email,
		// ImageURL: user.ImageURL,
	}
}
