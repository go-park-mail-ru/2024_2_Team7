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

	event, err := h.EventService.GetEventByID(r.Context(), &pb.GetEventByIDRequest{ID: int32(id)})
	if err != nil {
		if err != nil {
			st, ok := grpcStatus.FromError(err)
			if ok {
				switch st.Code() {
				case grpcCodes.NotFound:
					utils.WriteResponse(w, http.StatusConflict, httpErrors.ErrEventNotFound)
					return
				}
			}

			h.logger.Error(r.Context(), "get event by id", err)
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
			return
		}
	}

	resp := eventToEventResponse(event)
	utils.WriteResponse(w, http.StatusOK, resp)
}
