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

// @Summary Удаление события
// @Description Удаляет существующее событие
// @Tags events
// @Produce  json
// @Success 204
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Failure 403 {object} httpErrors.HttpError "Access Denied"
// @Failure 404 {object} httpErrors.HttpError "Event Not Found"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/{id} [delete]
func (h EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())

	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorID := session.UserID
	err = h.service.DeleteEvent(r.Context(), id, authorID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			utils.WriteResponse(w, http.StatusNotFound, httpErrors.ErrEventNotFound)
		case errors.Is(err, models.ErrAccessDenied):
			utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrAccessDenied)
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}
}
