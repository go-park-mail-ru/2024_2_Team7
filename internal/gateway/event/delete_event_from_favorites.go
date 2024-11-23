package events

import (
	"net/http"
	"strconv"

	pb "kudago/internal/event/api"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"

	"github.com/gorilla/mux"
)

// @Summary Удаление события из избранного
// @Description Удаляет событие из списка избранного
// @Tags events
// @Produce  json
// @Success 200
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/favorites/{id} [delete]
func (h EventHandler) DeleteEventFromFavorites(w http.ResponseWriter, r *http.Request) {
	session, ok := utils.GetSessionFromContext(r.Context())
	if !ok {
		utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidID)
		return
	}

	newFavorite := &pb.FavoriteEvent{
		UserID:  int32(session.UserID),
		EventID: int32(id),
	}

	_, err = h.EventService.DeleteEventFromFavorites(r.Context(), newFavorite)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
	}
	w.WriteHeader(http.StatusOK)
}
