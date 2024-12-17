package handlers

import (
	"net/http"

	httpErrors "kudago/internal/gateway/errors"
	pb "kudago/internal/user/api"

	"kudago/internal/gateway/utils"
)

// @Summary Получение избранных событий
// @Description Возвращает избранные события
// @Tags events
// @Produce  json
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество событий на странице (по умолчанию 30)"
// @Success 200 {object} GetEventsResponse
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/favorites [get]
func (h *UserHandlers) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	req := &pb.GetSubscribersRequest{
		ID: int32(session.UserID),
	}

	users, err := h.UserService.GetSubscribers(r.Context(), req)
	if err != nil {
		h.logger.Error(r.Context(), "getFavorites", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := writeUsersResponse(users.Users, len(users.Users))

	utils.WriteResponse(w, http.StatusOK, resp)
}
