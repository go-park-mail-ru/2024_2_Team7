package events

import (
	"net/http"
	"strconv"

	pb "kudago/internal/event/api"
	httpErrors "kudago/internal/http/errors"
	"kudago/internal/gateway/utils"

	"github.com/gorilla/mux"
)

// @Summary Добавление события в изсбранное
// @Description Добавить событие в избранное
// @Tags events
// @Produce  json
// @Success 200
// @Failure 401 {object} httpErrors.HttpError "Unauthorized"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/favorites/{id} [post]
func (h EventHandler) AddEventToFavorites(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.EventService.AddEventToFavorites(r.Context(), newFavorite)
	if err != nil {
		utils.WriteResponse(w, http.StatusConflict, err)
	}
	w.WriteHeader(http.StatusOK)
}
