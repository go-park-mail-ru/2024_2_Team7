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

// @Summary Получение событий по категори
// @Description Возвращает события по ID категории
// @Tags events
// @Produce  json
// @Success 200 {object} GetEventsResponse
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/categories/{category} [get]
func (h EventHandler) GetEventsByCategory(w http.ResponseWriter, r *http.Request) {
	paginationParams := GetPaginationParams(r)
	vars := mux.Vars(r)
	category := vars["category"]
	categoryID, err := strconv.Atoi(category)
	if err != nil {
		h.logger.Error(r.Context(), "getEventsByCategory", err)
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidCategory)
		return
	}

	req := &pb.GetEventsByCategoryRequest{
		CategoryID: int32(categoryID),
		Params:     paginationParams,
	}

	events, err := h.EventService.GetEventsByCategory(r.Context(), req)
	if err != nil {
		if err != nil {
			st, ok := grpcStatus.FromError(err)
			if ok {
				switch st.Code() {
				case grpcCodes.InvalidArgument:
					utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidCategory)
					return
				}
			}

			h.logger.Error(r.Context(), "get events by category", err)
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
			return
		}
	}

	resp := writeEventsResponse(events.Events, int(paginationParams.Limit))

	utils.WriteResponse(w, http.StatusOK, resp)
}
