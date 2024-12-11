package events

import (
	"net/http"
	"strconv"

	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	pb "kudago/internal/event/api"
	httpErrors "kudago/internal/gateway/errors"
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
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrEventNotFound)
				return
			case grpcCodes.AlreadyExists:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrEventAlreadyAddedToFavorites)
				return
			}
		}

		h.logger.Error(r.Context(), "add event to favorites", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}

	w.WriteHeader(http.StatusOK)
}
