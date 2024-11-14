package events

import (
	"net/http"
	"strconv"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newFavorite := models.FavoriteEvent{
		UserID:  session.UserID,
		EventID: id,
	}

	err = h.service.AddEventToFavorites(r.Context(), newFavorite)
	if err != nil {
		utils.WriteResponse(w, http.StatusConflict, err)
	}
	w.WriteHeader(http.StatusOK)
}