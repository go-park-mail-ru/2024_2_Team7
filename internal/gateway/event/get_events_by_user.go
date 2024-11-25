package events

import (
	"net/http"
	"strconv"

	pb "kudago/internal/event/api"
	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"

	"github.com/gorilla/mux"
)

// @Summary Получение событий пользователя
// @Description Возвращает события пользователя
// @Tags events
// @Produce  json
// @Success 200 {object} GetEventsResponse
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/user/{id} [get]
func (h EventHandler) GetEventsByUser(w http.ResponseWriter, r *http.Request) {
	paginationParams := GetPaginationParams(r)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidID)
		return
	}

	req := &pb.GetEventsByUserRequest{
		UserID: int32(id),
		Params: paginationParams,
	}

	events, err := h.EventService.GetEventsByUser(r.Context(), req)
	if err != nil {
		h.logger.Error(r.Context(), "get events by user", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := writeEventsResponse(events.Events, int(paginationParams.Limit))

	utils.WriteResponse(w, http.StatusOK, resp)
}
