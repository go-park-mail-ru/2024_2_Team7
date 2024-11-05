package events

import (
	"encoding/json"
	"net/http"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"
)

type SearchRequest struct {
	Query      string   `json:"query"`
	EventStart string   `json:"event_start"`
	EventEnd   string   `json:"event_end"`
	Tags       []string `json:"tags"`
	CategoryID int      `json:"category_id"`
}

// @Summary Поиск событий
// @Description Поиск событий по ключевым словам, датам, тегам и категории
// @Tags events
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество событий на странице (по умолчанию 30)"
// @Param SearchRequest body SearchRequest false "Фильтры для поиска событий"
// @Success 200 {object} GetEventsResponse "Список событий"
// @Failure 400 {object} httpErrors.HttpError "Invalid Data"
// @Failure 500 {object} httpErrors.HttpError "Internal Server Error"
// @Router /events [get]
func (h EventHandler) SearchEvents(w http.ResponseWriter, r *http.Request) {
	paginationParams := utils.GetPaginationParams(r)

	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(r.Context(), "search", err)
		utils.WriteResponse(w, http.StatusBadRequest, httpErrors.ErrInvalidData)
		return
	}

	params := models.SearchParams{
		Query:      req.Query,
		EventStart: req.EventStart,
		EventEnd:   req.EventEnd,
		Tags:       req.Tags,
		Category:   req.CategoryID,
	}

	events, err := h.service.SearchEvents(r.Context(), params, paginationParams)
	if err != nil {
		h.logger.Error(r.Context(), "search", err)
		utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrInternal)
		return
	}
	resp := writeEventsResponse(events, paginationParams.Limit)

	utils.WriteResponse(w, http.StatusOK, resp)
}
