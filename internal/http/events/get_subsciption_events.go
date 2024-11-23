package events

import (
	"net/http"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
)

// @Summary Получение событий по подпискам пользователя
// @Description Возвращает события пользователя
// @Tags events
// @Produce  json
// @Success 200 {object} GetEventsResponse
// @Failure 403 {object} httpErrors.HttpError "Status forbidden"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/subscription [get]
func (h *EventHandler) GetSubscriptionEvents(w http.ResponseWriter, r *http.Request) {
	paginationParams := utils.GetPaginationParams(r)

	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	events, err := h.getter.GetSubscriptionEvents(r.Context(), session.UserID, paginationParams)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	resp := writeEventsResponse(events, paginationParams.Limit)

	utils.WriteResponse(w, http.StatusOK, resp)
}
