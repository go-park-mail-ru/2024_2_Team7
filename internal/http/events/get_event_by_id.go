package events

import (
	"errors"
	"net/http"
	"strconv"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/gorilla/mux"
)

// @Summary Получение события по ID
// @Description Возвращает информацию о событии по его идентификатору
// @Tags events
// @Produce  json
// @Success 200 {object} EventResponse
// @Failure 404 {object} httpErrors.HttpError "Event Not Found"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/{id} [get]
func (h EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidID)
		return
	}

	event, err := h.getter.GetEventByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			w.WriteHeader(http.StatusNoContent)
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}
	resp := eventToEventResponse(event)
	utils.WriteResponse(w, http.StatusOK, resp)
}
