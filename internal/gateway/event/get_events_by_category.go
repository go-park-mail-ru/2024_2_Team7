package events

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	pb "kudago/internal/event/api"
	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"
)

// @Summary Получение событий по категори
// @Description Возвращает события по ID категории
// @Tags events
// @Produce  json
// @Success 200 {object} GetEventsResponse
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events/categories/{category} [get]
func (h EventHandler) GetEventsByCategory(w http.ResponseWriter, r *http.Request) {
	paginationParams := utils.GetPaginationParams(r)
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
		Params:     &pb.PaginationParams{},
	}

	events, err := h.EventService.GetEventsByCategory(r.Context(), req)
	if err != nil {
		h.logger.Error(r.Context(), "getEventsByCategory", err)
		switch err {
		case models.ErrInvalidCategory:
			utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidCategory)

		///TODO пока оставлю так, когда будет более четкая бд и ошибки для обработки, поправлю
		default:
			utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		}
		return
	}

	resp := writeEventsResponse(events.Events, paginationParams.Limit)

	utils.WriteResponse(w, http.StatusOK, resp)
}
