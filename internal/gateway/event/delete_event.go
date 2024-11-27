package events

import (
	"net/http"
	"strconv"

	pb "kudago/internal/event/api"
	httpErrors "kudago/internal/gateway/errors"
	"kudago/internal/gateway/utils"

	"github.com/gorilla/mux"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
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
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidID)
		return
	}

	authorID := session.UserID

	req := &pb.DeleteEventRequest{
		AuthorID: int32(authorID),
		EventID:  int32(id),
	}
	_, err = h.EventService.DeleteEvent(r.Context(), req)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if ok {
			switch st.Code() {
			case grpcCodes.NotFound:
				utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrEventNotFound)
				return
			}
		}

		h.logger.Error(r.Context(), "delete event", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}
}
