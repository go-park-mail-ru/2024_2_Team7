package handlers

import (
	"net/http"
	"strconv"

	pb "kudago/internal/user/api"

	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"

	"github.com/gorilla/mux"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
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
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrUserNotFound)
				return
			}
		}

		h.logger.Error(r.Context(), "profile", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
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
		ImageURL: user.AvatarUrl,
	}
}
