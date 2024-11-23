package events

import (
	"net/http"

	pb "kudago/internal/event/api"
	httpErrors "kudago/internal/http/errors"

	"kudago/internal/http/utils"
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
func (h EventHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	paginationParams := GetPaginationParams(r)

	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	req := &pb.GetFavoritesRequest{
		UserID: int32(session.UserID),
		Params: paginationParams,
	}

	events, err := h.EventService.GetFavorites(r.Context(), req)
	if err != nil {
		h.logger.Error(r.Context(), "getFavorites", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := writeEventsResponse(events.Events, int(paginationParams.Limit))

	utils.WriteResponse(w, http.StatusOK, resp)
}
